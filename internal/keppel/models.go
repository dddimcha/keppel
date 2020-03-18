/*******************************************************************************
*
* Copyright 2018 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package keppel

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/opencontainers/go-digest"
	gorp "gopkg.in/gorp.v2"
)

//Account contains a record from the `accounts` table.
type Account struct {
	Name         string `db:"name"`
	AuthTenantID string `db:"auth_tenant_id"`

	//UpstreamPeerHostName is set if and only if the "on_first_use" replication strategy is used.
	UpstreamPeerHostName string `db:"upstream_peer_hostname"`
	//RequiredLabels is a comma-separated list of labels that must be present on
	//all image manifests in this account.
	RequiredLabels string `db:"required_labels"`

	BlobsSweepedAt   *time.Time `db:"blobs_sweeped_at"`   //see tasks.SweepBlobsInNextAccount
	StorageSweepedAt *time.Time `db:"storage_sweeped_at"` //see tasks.SweepStorageInNextAccount
}

//SwiftContainerName returns the name of the Swift container backing this
//Keppel account.
func (a Account) SwiftContainerName() string {
	return "keppel-" + a.Name
}

//PostgresDatabaseName returns the name of the Postgres database which contains this
//Keppel account's metadata.
func (a Account) PostgresDatabaseName() string {
	return "keppel_" + strings.Replace(a.Name, "-", "_", -1)
}

//FindAccount works similar to db.SelectOne(), but returns nil instead of
//sql.ErrNoRows if no account exists with this name.
func FindAccount(db gorp.SqlExecutor, name string) (*Account, error) {
	var account Account
	err := db.SelectOne(&account,
		"SELECT * FROM accounts WHERE name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &account, err
}

////////////////////////////////////////////////////////////////////////////////

//RBACPolicy contains a record from the `rbac_policies` table.
type RBACPolicy struct {
	AccountName        string `db:"account_name"`
	RepositoryPattern  string `db:"match_repository"`
	UserNamePattern    string `db:"match_username"`
	CanPullAnonymously bool   `db:"can_anon_pull"`
	CanPull            bool   `db:"can_pull"`
	CanPush            bool   `db:"can_push"`
	CanDelete          bool   `db:"can_delete"`
}

//Matches evaluates the regexes in this policy.
func (r RBACPolicy) Matches(repoName, userName string) bool {
	if r.RepositoryPattern != "" {
		rx, err := regexp.Compile(fmt.Sprintf(`^%s/%s$`,
			regexp.QuoteMeta(r.AccountName),
			r.RepositoryPattern,
		))
		if err != nil || !rx.MatchString(repoName) {
			return false
		}
	}

	if r.UserNamePattern != "" {
		rx, err := regexp.Compile(fmt.Sprintf(`^%s$`, r.UserNamePattern))
		if err != nil || !rx.MatchString(userName) {
			return false
		}
	}

	return true
}

////////////////////////////////////////////////////////////////////////////////

//Blob contains a record from the `blobs` table.
//
//In the `blobs` table, blobs are only bound to an account. This makes
//cross-repo blob mounts cheap and easy to implement. The actual connection to
//repos is in the `blob_mounts` table.
//
//StorageID is used to construct the filename (or equivalent) for this blob
//in the StorageDriver. We cannot use the digest for this since the StorageID
//needs to be chosen at the start of the blob upload, when the digest is not
//known yet.
type Blob struct {
	ID                     int64      `db:"id"`
	AccountName            string     `db:"account_name"`
	Digest                 string     `db:"digest"`
	SizeBytes              uint64     `db:"size_bytes"`
	StorageID              string     `db:"storage_id"`
	PushedAt               time.Time  `db:"pushed_at"`
	ValidatedAt            time.Time  `db:"validated_at"` //see tasks.ValidateNextBlob
	ValidationErrorMessage string     `db:"validation_error_message"`
	MarkedForDeletionAt    *time.Time `db:"marked_for_deletion_at"` //see tasks.SweepBlobsInNextAccount
}

const blobGetQueryByRepoName = `
	SELECT b.*
	  FROM blobs b
	  JOIN blob_mounts bm ON b.id = bm.blob_id
	  JOIN repos r ON bm.repo_id = r.id
	 WHERE b.account_name = $1 AND b.digest = $2
	   AND r.account_name = $1 AND r.name = $3
`

const blobGetQueryByRepoID = `
	SELECT b.*
	  FROM blobs b
	  JOIN blob_mounts bm ON b.id = bm.blob_id
	 WHERE b.account_name = $1 AND b.digest = $2 AND bm.repo_id = $3
`

const blobGetQueryByAccountName = `
	SELECT * FROM blobs WHERE account_name = $1 AND digest = $2
`

//FindBlobByRepositoryName is a convenience wrapper around db.SelectOne(). If
//the blob in question does not exist, sql.ErrNoRows is returned.
func FindBlobByRepositoryName(db gorp.SqlExecutor, blobDigest digest.Digest, repoName string, account Account) (*Blob, error) {
	var blob Blob
	err := db.SelectOne(&blob, blobGetQueryByRepoName, account.Name, blobDigest.String(), repoName)
	return &blob, err
}

//FindBlobByRepository is a convenience wrapper around db.SelectOne(). If
//the blob in question does not exist, sql.ErrNoRows is returned.
func FindBlobByRepository(db gorp.SqlExecutor, blobDigest digest.Digest, repo Repository, account Account) (*Blob, error) {
	var blob Blob
	err := db.SelectOne(&blob, blobGetQueryByRepoID, account.Name, blobDigest.String(), repo.ID)
	return &blob, err
}

//FindBlobByAccountName is a convenience wrapper around db.SelectOne(). If the
//blob in question does not exist, sql.ErrNoRows is returned.
func FindBlobByAccountName(db gorp.SqlExecutor, blobDigest digest.Digest, account Account) (*Blob, error) {
	var blob Blob
	err := db.SelectOne(&blob, blobGetQueryByAccountName, account.Name, blobDigest.String())
	return &blob, err
}

//MountBlobIntoRepo creates an entry in the blob_mounts database table.
func MountBlobIntoRepo(db gorp.SqlExecutor, blob Blob, repo Repository) error {
	_, err := db.Exec(
		`INSERT INTO blob_mounts (blob_id, repo_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		blob.ID, repo.ID,
	)
	return err
}

//Upload contains a record from the `uploads` table.
//
//Digest contains the SHA256 digest of everything that has been uploaded so
//far. This is used to validate that we're resuming at the right position in
//the next PUT/PATCH.
type Upload struct {
	RepositoryID int64     `db:"repo_id"`
	UUID         string    `db:"uuid"`
	StorageID    string    `db:"storage_id"`
	SizeBytes    uint64    `db:"size_bytes"`
	Digest       string    `db:"digest"`
	NumChunks    uint32    `db:"num_chunks"`
	UpdatedAt    time.Time `db:"updated_at"`
}

const uploadGetQueryByRepoID = `
	SELECT u.* FROM uploads u WHERE u.uuid = $1 AND repo_id = $2
`

//FindUploadByRepository is a convenience wrapper around db.SelectOne(). If
//the upload in question does not exist, sql.ErrNoRows is returned.
func FindUploadByRepository(db gorp.SqlExecutor, uuid string, repo Repository) (*Upload, error) {
	var upload Upload
	err := db.SelectOne(&upload, uploadGetQueryByRepoID, uuid, repo.ID)
	return &upload, err
}

////////////////////////////////////////////////////////////////////////////////

//Repository contains a record from the `repos` table.
type Repository struct {
	ID                  int64      `db:"id"`
	AccountName         string     `db:"account_name"`
	Name                string     `db:"name"`
	BlobMountsSweepedAt *time.Time `db:"blob_mounts_sweeped_at"` //see tasks.SweepBlobMountsInNextRepo
	ManifestsSyncedAt   *time.Time `db:"manifests_synced_at"`    //see tasks.SyncManifestsInNextRepo (only set for replica accounts)
}

//FindOrCreateRepository works similar to db.SelectOne(), but autovivifies a
//Repository record when none exists yet.
func FindOrCreateRepository(db gorp.SqlExecutor, name string, account Account) (*Repository, error) {
	repo, err := FindRepository(db, name, account)
	if err == sql.ErrNoRows {
		repo = &Repository{
			AccountName: account.Name,
			Name:        name,
		}
		err = db.Insert(repo)
	}
	return repo, err
}

//FindRepository is a convenience wrapper around db.SelectOne(). If the
//repository in question does not exist, sql.ErrNoRows is returned.
func FindRepository(db gorp.SqlExecutor, name string, account Account) (*Repository, error) {
	var repo Repository
	err := db.SelectOne(&repo,
		"SELECT * FROM repos WHERE account_name = $1 AND name = $2", account.Name, name)
	return &repo, err
}

//FullName prepends the account name to the repository name.
func (r Repository) FullName() string {
	return r.AccountName + `/` + r.Name
}

////////////////////////////////////////////////////////////////////////////////

//Manifest contains a record from the `manifests` table.
type Manifest struct {
	RepositoryID           int64     `db:"repo_id"`
	Digest                 string    `db:"digest"`
	MediaType              string    `db:"media_type"`
	SizeBytes              uint64    `db:"size_bytes"`
	PushedAt               time.Time `db:"pushed_at"`
	ValidatedAt            time.Time `db:"validated_at"` //see tasks.ValidateNextManifest
	ValidationErrorMessage string    `db:"validation_error_message"`
}

//FindManifest is a convenience wrapper around db.SelectOne(). If the
//manifest in question does not exist, sql.ErrNoRows is returned.
func FindManifest(db gorp.SqlExecutor, repo Repository, digest string) (*Manifest, error) {
	var manifest Manifest
	err := db.SelectOne(&manifest,
		"SELECT * FROM manifests WHERE repo_id = $1 AND digest = $2", repo.ID, digest)
	return &manifest, err
}

//Tag contains a record from the `tags` table.
type Tag struct {
	RepositoryID int64     `db:"repo_id"`
	Name         string    `db:"name"`
	Digest       string    `db:"digest"`
	PushedAt     time.Time `db:"pushed_at"`
}

////////////////////////////////////////////////////////////////////////////////

//Quotas contains a record from the `quotas` table.
type Quotas struct {
	AuthTenantID  string `db:"auth_tenant_id"`
	ManifestCount uint64 `db:"manifests"`
}

//FindQuotas works similar to db.SelectOne(), but returns nil instead of
//sql.ErrNoRows if no quota set exists for this auth tenant.
func FindQuotas(db gorp.SqlExecutor, authTenantID string) (*Quotas, error) {
	var quotas Quotas
	err := db.SelectOne(&quotas,
		"SELECT * FROM quotas WHERE auth_tenant_id = $1", authTenantID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &quotas, err
}

//DefaultQuotas creates a new Quotas instance with the default quotas.
func DefaultQuotas(authTenantID string) *Quotas {
	//Right now, the default quota is always 0. The value of having this function
	//is to ensure that we only need to change this place if this ever changes.
	return &Quotas{
		AuthTenantID:  authTenantID,
		ManifestCount: 0,
	}
}

var manifestUsageQuery = `
	SELECT COUNT(m.digest)
	  FROM manifests m
	  JOIN repos r ON m.repo_id = r.id
	  JOIN accounts a ON a.name = r.account_name
	 WHERE a.auth_tenant_id = $1
`

//GetManifestUsage returns how many manifests currently exist in repos in
//accounts connected to this quota set's auth tenant.
func (q Quotas) GetManifestUsage(db gorp.SqlExecutor) (uint64, error) {
	manifestCount, err := db.SelectInt(manifestUsageQuery, q.AuthTenantID)
	return uint64(manifestCount), err
}

////////////////////////////////////////////////////////////////////////////////

//Peer contains a record from the `peers` table.
type Peer struct {
	HostName string `db:"hostname"`

	//OurPassword is what we use to log in at the peer.
	OurPassword string `db:"our_password"`

	//TheirCurrentPasswordHash and TheirPreviousPasswordHash is what the peer
	//uses to log in with us. Passwords are rotated hourly. We allow access with
	//the current *and* the previous password to avoid a race where we enter the
	//new password in the database and then reject authentication attempts from
	//the peer before we told them about the new password.
	TheirCurrentPasswordHash  string `db:"their_current_password_hash"`
	TheirPreviousPasswordHash string `db:"their_previous_password_hash"`

	//LastPeeredAt is when we last issued a new password for this peer.
	LastPeeredAt *time.Time `db:"last_peered_at"` //see tasks.IssueNewPasswordForPeer
}

////////////////////////////////////////////////////////////////////////////////

//PendingBlob contains a record from the `pending_blobs` table.
type PendingBlob struct {
	AccountName  string        `db:"account_name"`
	Digest       string        `db:"digest"`
	Reason       PendingReason `db:"reason"`
	PendingSince time.Time     `db:"since"`
}

//PendingReason is an enum that explains why a blob is pending.
type PendingReason string

const (
	//PendingBecauseOfReplication is when a blob is pending because
	//it is currently being replicated from an upstream registry.
	PendingBecauseOfReplication PendingReason = "replication"
)

////////////////////////////////////////////////////////////////////////////////

//UnknownBlob contains a record from the `unknown_blobs` table.
//This is only used by tasks.SweepStorageInNextAccount().
type UnknownBlob struct {
	AccountName         string    `db:"account_name"`
	StorageID           string    `db:"storage_id"`
	MarkedForDeletionAt time.Time `db:"marked_for_deletion_at"`
}

//UnknownManifest contains a record from the `unknown_manifests` table.
//This is only used by tasks.SweepStorageInNextAccount().
//
//NOTE: We don't use repository IDs here because unknown manifests may exist in
//repositories that are also not known to the database.
type UnknownManifest struct {
	AccountName         string    `db:"account_name"`
	RepositoryName      string    `db:"repo_name"`
	Digest              string    `db:"digest"`
	MarkedForDeletionAt time.Time `db:"marked_for_deletion_at"`
}

////////////////////////////////////////////////////////////////////////////////

func initModels(db *gorp.DbMap) {
	db.AddTableWithName(Account{}, "accounts").SetKeys(false, "name")
	db.AddTableWithName(RBACPolicy{}, "rbac_policies").SetKeys(false, "account_name", "match_repository", "match_username")
	db.AddTableWithName(Blob{}, "blobs").SetKeys(true, "id")
	db.AddTableWithName(Upload{}, "uploads").SetKeys(false, "repo_id", "uuid")
	db.AddTableWithName(Repository{}, "repos").SetKeys(true, "id")
	db.AddTableWithName(Manifest{}, "manifests").SetKeys(false, "repo_id", "digest")
	db.AddTableWithName(Tag{}, "tags").SetKeys(false, "repo_id", "name")
	db.AddTableWithName(Quotas{}, "quotas").SetKeys(false, "auth_tenant_id")
	db.AddTableWithName(Peer{}, "peers").SetKeys(false, "hostname")
	db.AddTableWithName(PendingBlob{}, "pending_blobs").SetKeys(false, "account_name", "digest")
	db.AddTableWithName(UnknownBlob{}, "unknown_blobs").SetKeys(false, "account_name", "storage_id")
	db.AddTableWithName(UnknownManifest{}, "unknown_manifests").SetKeys(false, "account_name", "repo_name", "digest")
}
