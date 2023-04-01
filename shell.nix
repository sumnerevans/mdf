{ pkgs ? import <nixpkgs> { } }: with pkgs;
pkgs.mkShell {
  buildInputs = [
    libnotify
    go
    gotools
    gopls
    yq
  ];
}
