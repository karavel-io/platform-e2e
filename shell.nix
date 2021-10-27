let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
    buildInputs = with pkgs; [
        nodejs-14_x
        kubectl
        kustomize
        nodePackages.pnpm
        awscli2
        doctl
        unstable.terraform
    ];
}
