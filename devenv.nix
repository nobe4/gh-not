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
    # golangci-lint moves fast, and I like to have its latest rules.
    pkgs-unstable.golangci-lint
  ];
}
