{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
        };

        go-migrate-postgres = pkgs.go-migrate.overrideAttrs (_old: {
          tags = ["postgres"];
        });
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go

            # Go development tools.
            go-task
            sqlc
            go-migrate-postgres
            golangci-lint
            govulncheck

            # Protobuf / ConnectRPC code generation tools.
            buf
            protoc-gen-go
            protoc-gen-connect-go

            # Load test tools.
            vegeta
          ];

          shellHook = ''
            echo "go: $(go version)"
            echo "task: $(task --version)"
            echo "sqlc: $(sqlc version)"
            echo "migrate: $(migrate -version 2>&1 | head -n 1)"
            echo "golangci-lint: $(golangci-lint --version | head -n 1)"
            echo "govulncheck: $(govulncheck -version 2>&1 | head -n 1)"
            echo "buf: $(buf --version)"
            echo "protoc-gen-go: $(protoc-gen-go --version)"
            echo "protoc-gen-connect-go: $(protoc-gen-connect-go --version)"
            echo "vegeta: $(vegeta --version 2>&1 | head -n 1)"
            echo ""
          '';
        };
      }
    );
}
