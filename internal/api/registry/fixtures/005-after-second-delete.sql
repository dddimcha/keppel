INSERT INTO accounts (name, auth_tenant_id, registry_http_secret, upstream_peer_hostname) VALUES ('test1', 'test1authtenant', 'topsecret', '');

INSERT INTO blob_mounts (blob_id, repo_id) VALUES (1, 1);
INSERT INTO blob_mounts (blob_id, repo_id) VALUES (2, 1);

INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at) VALUES (1, 'test1', 'sha256:7575de20fdeacfb9a529c26f03c5beab401bb985721b923bba6b34fe4fce5f9c', 1486, '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b', 3);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at) VALUES (2, 'test1', 'sha256:27329396f45bf916fc3c4af52592ca88235f34bb7b4475bc52ed452058d8d160', 1518, 'd4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35', 4);

INSERT INTO quotas (auth_tenant_id, manifests) VALUES ('test1authtenant', 100);

INSERT INTO repos (id, account_name, name) VALUES (1, 'test1', 'foo');
