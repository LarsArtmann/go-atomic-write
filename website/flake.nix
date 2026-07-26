{
  description = "go-atomic-write website — Astro + Starlight";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    systems.url = "github:nix-systems/default";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        { config, pkgs, ... }:
        let
          mkApp =
            name: description: runtimeInputs: text:
            {
              type = "app";
              program = "${
                pkgs.writeShellApplication {
                  inherit name runtimeInputs text;
                }
              }/bin/${name}";
              meta = {
                inherit description;
                mainProgram = name;
              };
            };
        in
        {
          apps = {
            dev = mkApp "dev" "Start the Astro development server" [ pkgs.nodejs_24 ] "npm run dev";
            build = mkApp "build" "Build the website for production" [ pkgs.nodejs_24 ] "npm run build";
            preview = mkApp "preview" "Preview the production build locally" [ pkgs.nodejs_24 ] "npm run preview";
            deploy = mkApp "deploy" "Build and deploy the website to Firebase Hosting" [
              pkgs.nodejs_24
              pkgs.firebase-tools
            ] ''
              npm run build
              firebase deploy --only hosting
            '';
          };

          devShells.default = pkgs.mkShellNoCC {
            packages = builtins.attrValues {
              inherit (pkgs) nodejs_24 firebase-tools;
            };
          };

          treefmt.programs.nixfmt.enable = true;

          checks.format = config.treefmt.build.check self;
        };
    };
}
