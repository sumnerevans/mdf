{
  description = "Sumner's mutt display filter";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in rec {
        packages.mdf = pkgs.buildGoModule {
          pname = "mdf";
          version = "unstable-2023-12-04";
          src = self;

          vendorHash = "sha256-CuO80I648lrHpaJ+T4yWwfmX6J1wDpkr2mVGHMybt0A=";
        };
        defaultPackage = packages.mdf;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ libnotify go gotools gopls yq ];
        };
      }));
}
