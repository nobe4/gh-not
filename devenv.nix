{
  pkgs,
  inputs,
  ...
}:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  packages = [
    pkgs.git
    pkgs.go
    pkgs.golangci-lint-langserver
    pkgs-unstable.golangci-lint
  ];
}
