{
  description = "Sumner's mutt display filter";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in {
        packages = {
          mdf = pkgs.buildGoModule {
            pname = "mdf";
            version = "unstable-2025-10-11";
            src = self;

            vendorHash = "sha256-zIZMPA/vYD+SU1+//KKcLooCPL6eWxTVkSSAIG7cjtk=";
          };
          default = self.outputs.packages.${system}.mdf;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ libnotify go gotools gopls yq ];
        };
      }));
}
