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

          vendorHash = "sha256-PcXa+5cvprr9h0RaGvlSG5GtNYT7A8pn3sD0neisHec=";
        };
        defaultPackage = packages.mdf;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ libnotify go gotools gopls yq ];
        };
      }));
}
