{
  description = "A CLI for managing Linode infrastructure";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "linodectl";
          version = "0.0.1";
          src = ./.;
          vendorHash = "sha256-rphnh6xNnZTwetMrBYNzaCkeXGCLXyrBrAvp4wDIEAY=";
          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "A CLI for managing Linode infrastructure";
            homepage = "https://github.com/0xch4z/linodectl";
            license = licenses.mit;
            mainProgram = "linodectl";
          };
        };
      }
    )
    // {
      overlays.default = final: prev: {
        linodectl = self.packages.${prev.system}.default;
      };
    };
}
