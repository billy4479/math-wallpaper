This is a collection of scripts I've put together to generate an image full of math formulas.

Yes, I'm a nerd.

This code absolutely sucks and assumes makes a lot of assumptions:

- You should have a bunch of latex in `../latex/`
- Your latex should use `\begin{equation}` and not `$$`or `\[`
- You have a [latex-build](https://github.com/billy4479/latex-build) installed
  - You can use anything that builds latex, but this is the recommended way since it builds in parallel
  - With 8 workers it took me 124s to compile 675 files
- You know how to tweak the values in `./render_image.go` to customize them to your liking

If all of that is true you can just™ run

```sh
go build
./math-wallpaper extract
latex-build
./math-wallpaper render
```
