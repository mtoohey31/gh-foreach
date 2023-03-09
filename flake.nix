{
  description = "gh-foreach";
  inputs = {
    utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  };
  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem
      (system:
        with import nixpkgs { inherit system; }; {
          packages.default = buildGo118Module rec {
            name = "gh-foreach";
            pname = name;
            src = ./.;
            vendorSha256 = "DZDydQQGwwDyX0yGh/8AmR6aBzkJj3HbNP6JC5dPCaE=";
          };

          devShells.default = mkShell {
            nativeBuildInputs = [ go_1_18 gopls ];
          };
        }) // {
      overlays.default = (final: _: {
        gh-foreach = self.packages."${final.system}".default;
      });
    };
}
