INSERT INTO accounts (name, auth_tenant_id, upstream_peer_hostname, required_labels, metadata_json, next_blob_sweep_at, next_storage_sweep_at, next_federation_announcement_at, in_maintenance) VALUES ('test1', 'test1authtenant', 'registry.example.org', '', '', NULL, NULL, NULL, TRUE);

INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (1, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (2, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (3, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (4, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (5, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (6, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (7, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (8, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, can_be_deleted_at) VALUES (9, 1, NULL);

INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (1, 'test1', 'sha256:aa9bf3251d0754379f51f5a7ca15835c504a1b1e8fa6663463c22dfc9ae527b8', 1048576, 'aa9bf3251d0754379f51f5a7ca15835c504a1b1e8fa6663463c22dfc9ae527b8', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (2, 'test1', 'sha256:5dfe56df10488d9c1f75414293fe474eaa7fca63acdee374bc371250f1a372b4', 1048576, '5dfe56df10488d9c1f75414293fe474eaa7fca63acdee374bc371250f1a372b4', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (3, 'test1', 'sha256:db8bc83bac14839cc7d46e346e417a9ecd115732205c79311c5e191dd156ed18', 1412, 'db8bc83bac14839cc7d46e346e417a9ecd115732205c79311c5e191dd156ed18', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (4, 'test1', 'sha256:4b2ac57bae1cdb45d45edc74ff4b2c28f60434ff3ce128b4e0bd4034528046eb', 1048576, '4b2ac57bae1cdb45d45edc74ff4b2c28f60434ff3ce128b4e0bd4034528046eb', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (5, 'test1', 'sha256:939fc3e1f9c5e776c28b29b807fb650274ae0eb8bc04889ecbddb38100b977e2', 1048576, '939fc3e1f9c5e776c28b29b807fb650274ae0eb8bc04889ecbddb38100b977e2', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (6, 'test1', 'sha256:7b4640847d186b4d5e10ba0335f09d2a06e1d55b5a958919a55b379bd4a018f1', 1412, '7b4640847d186b4d5e10ba0335f09d2a06e1d55b5a958919a55b379bd4a018f1', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (7, 'test1', 'sha256:8e179b08fa229a8c2958d000ca68113cdd7860d38338f376c131c4a52e550140', 1048576, '8e179b08fa229a8c2958d000ca68113cdd7860d38338f376c131c4a52e550140', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (8, 'test1', 'sha256:6ea2d07e23765eb41a6577d42905eccf55b6697cb9dd2a544333766b65748ec5', 1048576, '6ea2d07e23765eb41a6577d42905eccf55b6697cb9dd2a544333766b65748ec5', 3600, 3600, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, can_be_deleted_at) VALUES (9, 'test1', 'sha256:5a7ef6200915bdc81598c5d700001f16dfe0f4d2ff3a62bfa89ebfb534d5bd55', 1412, '5a7ef6200915bdc81598c5d700001f16dfe0f4d2ff3a62bfa89ebfb534d5bd55', 3600, 3600, '', NULL);

INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:33ef30d47bd666b28f971cc3f07b00aca72d55865e79d6ca03a836647bb81f42', 1);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:33ef30d47bd666b28f971cc3f07b00aca72d55865e79d6ca03a836647bb81f42', 2);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:33ef30d47bd666b28f971cc3f07b00aca72d55865e79d6ca03a836647bb81f42', 3);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:5e81574b620cb58bd6db9ba7cfab0eefc6fdf36ec04a423af737a23c0b481d57', 7);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:5e81574b620cb58bd6db9ba7cfab0eefc6fdf36ec04a423af737a23c0b481d57', 8);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:5e81574b620cb58bd6db9ba7cfab0eefc6fdf36ec04a423af737a23c0b481d57', 9);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:622d6cd420f93b90fcf415b3254eb798eb1df51c7aecbf1ccc8331acea30079d', 4);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:622d6cd420f93b90fcf415b3254eb798eb1df51c7aecbf1ccc8331acea30079d', 5);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:622d6cd420f93b90fcf415b3254eb798eb1df51c7aecbf1ccc8331acea30079d', 6);

INSERT INTO manifest_manifest_refs (repo_id, parent_digest, child_digest) VALUES (1, 'sha256:cb117c93c9ebf1b0b7eb0be4f0abfee8ec773a895fa64dc7ea13b824a6b11506', 'sha256:33ef30d47bd666b28f971cc3f07b00aca72d55865e79d6ca03a836647bb81f42');
INSERT INTO manifest_manifest_refs (repo_id, parent_digest, child_digest) VALUES (1, 'sha256:cb117c93c9ebf1b0b7eb0be4f0abfee8ec773a895fa64dc7ea13b824a6b11506', 'sha256:622d6cd420f93b90fcf415b3254eb798eb1df51c7aecbf1ccc8331acea30079d');

INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at, validation_error_message, last_pulled_at) VALUES (1, 'sha256:33ef30d47bd666b28f971cc3f07b00aca72d55865e79d6ca03a836647bb81f42', 'application/vnd.docker.distribution.manifest.v2+json', 2099156, 3600, 3600, '', NULL);
INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at, validation_error_message, last_pulled_at) VALUES (1, 'sha256:5e81574b620cb58bd6db9ba7cfab0eefc6fdf36ec04a423af737a23c0b481d57', 'application/vnd.docker.distribution.manifest.v2+json', 2099156, 3600, 3600, '', NULL);
INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at, validation_error_message, last_pulled_at) VALUES (1, 'sha256:622d6cd420f93b90fcf415b3254eb798eb1df51c7aecbf1ccc8331acea30079d', 'application/vnd.docker.distribution.manifest.v2+json', 2099156, 3600, 3600, '', NULL);
INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at, validation_error_message, last_pulled_at) VALUES (1, 'sha256:cb117c93c9ebf1b0b7eb0be4f0abfee8ec773a895fa64dc7ea13b824a6b11506', 'application/vnd.docker.distribution.manifest.list.v2+json', 1715, 3600, 3600, '', NULL);

INSERT INTO peers (hostname, our_password, their_current_password_hash, their_previous_password_hash, last_peered_at) VALUES ('registry.example.org', 'a4cb6fae5b8bb91b0b993486937103dab05eca93', '', '', NULL);

INSERT INTO repos (id, account_name, name, next_blob_mount_sweep_at, next_manifest_sync_at) VALUES (1, 'test1', 'foo', NULL, 14400);
