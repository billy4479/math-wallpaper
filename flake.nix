{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
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

        texPkgs = pkgs.texliveSmall.withPackages (
          ps: with ps; [
            latexmk
            latexindent
            texcount

            # Direct dependencies
            dirtytalk
            tcolorbox
            tikzfill
            enumitem
            doclicense
            cancel
            minted
            svg
            physics
            rsfs
            lipsum
            cleveref
            thmtools
            siunitx
            subfig
            emoji
            doublestroke
            circuitikz
            pgfplots
            standalone

            # Indirect dependencies
            environ
            pdfcol
            xifthen
            ifmtarg
            xstring
            ccicons
            csquotes
            upquote
            transparent
            hyperxmp
            luacode
            luatex85
          ]
        );
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go

            texPkgs
          ];
        };
      }
    );
}
