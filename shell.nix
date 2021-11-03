let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
    buildInputs = with pkgs; [
        nodejs-14_x
        go
        go-junit-report
        sonobuoy
        kubectl
        kustomize
        nodePackages.pnpm
    ];
}
