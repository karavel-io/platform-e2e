let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
    buildInputs = with pkgs; [
        nodejs-14_x
        go
        sonobuoy
        kubectl
        kustomize
        nodePackages.pnpm
    ];
}
