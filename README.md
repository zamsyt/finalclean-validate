# TFC-validate (2023)

## Instructions v0.3.0

1. Save the Drawpile as a .ora file
2. With Go installed, in the repo root, run `go run . <command> [options]`

### `fullmerge`

`fullmerge <path-to-drawpile.ora>`

Convert the entire merged image to the palette, and save an image of the (corrected) pixels that are different from "BASE LAYER" to diff.png


### `split`

`split <my-image.png>`

Splits the (3000x2000) image to 125x125px sectors. Sector images are saved as `split/<major><minor>.png`

### `join`

`join [dir]`

Joins the images created by `split` back into `combined.png`
