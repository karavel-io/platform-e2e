let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
    buildInputs = with pkgs; [
        nodejs-14_x
        go_1_17
        kubectl
        kustomize
        nodePackages.pnpm
    ];
}
