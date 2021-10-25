package navigator

import (
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Some special "extension" values.
const (
	GenericFile  = "file"
	ClosedFolder = "folder-closed"
	OpenFolder   = "folder-open"
)

var fileTypes map[string]*unison.SvgPath

// InitFileTypes initializes the file type associations.
func InitFileTypes() {
	fileTypes = make(map[string]*unison.SvgPath)
	addFileType(ClosedFolder, "M464 128H272l-64-64H48C21.49 64 0 85.49 0 112v288c0 26.51 21.49 48 48 48h416c26.51 0 48-21.49 48-48V176c0-26.51-21.49-48-48-48z", 512, 512)
	addFileType(OpenFolder, "M572.694 292.093L500.27 416.248A63.997 63.997 0 0 1 444.989 448H45.025c-18.523 0-30.064-20.093-20.731-36.093l72.424-124.155A64 64 0 0 1 152 256h399.964c18.523 0 30.064 20.093 20.73 36.093zM152 224h328v-48c0-26.51-21.49-48-48-48H272l-64-64H48C21.49 64 0 85.49 0 112v278.046l69.077-118.418C86.214 242.25 117.989 224 152 224z", 576, 512)
	fileSvg := "M224 136V0H24C10.7 0 0 10.7 0 24v464c0 13.3 10.7 24 24 24h336c13.3 0 24-10.7 24-24V160H248c-13.2 0-24-10.8-24-24zm160-14.1v6.1H256V0h6.1c6.4 0 12.5 2.5 17 7l97.9 98c4.5 4.5 7 10.6 7 16.9z"
	addFileType(GenericFile, fileSvg, 384, 512)
	addFileType(".gcs", "M255 45.4c-24.5 0-47 11.8-63.9 33.4-16.9 21.5-27.1 52.6-27.1 86.5 0 36 12.1 67.5 31 89.5l13.5 15-19.6 4.6c-52.3 11.9-77.4 36.9-91.75 75.2-13.7 35.7-15.6 84.8-16.1 143.3H431c-.2-58.7-.5-109.3-13-145.5-13.4-39.4-37.9-64.3-94-75.4l-19.9-3.7 12.9-15.7c17.7-21.9 28.8-52.6 28.8-87.5 0-33.9-10.3-64.9-27.2-86.3-16.8-21.7-39.3-33.6-63.6-33.4z", 512, 512)
	addFileType(".gct", "M250.322 18.494c-25.06 3.26-47.158 32.267-47.158 69.346 0 20.453 7.06 38.57 17.502 51.166l10.123 12.213-15.59 2.932c-13.676 2.574-23.794 9.896-32.272 21.547-8.48 11.65-14.86 27.7-19.326 46.095-8.23 33.9-9.916 75.216-10.143 111.275h44.007l11.883 159.512h96.37l10.514-159.512h41.88c-.013-36.448-.353-78.316-7.81-112.48-4.042-18.524-10.176-34.575-18.777-46.12-8.6-11.543-19.21-18.81-34.482-21.18l-15.912-2.468 10.037-12.59c9.99-12.533 16.7-30.436 16.7-50.392 0-39.537-24.776-69.268-52.352-69.268-2.915 0-4.754-.135-5.196-.078zm178.608 1.078c-31.872-.534-61.166 26.473-71.084 63.49-4.575 17.073-4.83 35.29-.817 51.108-10.96 1.307-20.99 5.173-29.772 10.996 5.563 3.58 10.537 7.906 14.906 12.814 7.998-4.296 16.716-6.28 27.084-5.492l15.816 1.2-6.615-14.415c-5.86-12.764-7.33-33.55-2.554-51.377 8.122-30.308 31.484-49.75 52.75-49.61 1.416.008 2.825.104 4.22.29l.01.002c.263.037 1.817.567 4.44 1.27 23.73 6.36 38.404 37.853 29.168 72.324-4.66 17.392-15.965 34.567-27.02 42.73l-12.954 9.565 14.73 6.502c13.063 5.765 20.835 13.86 25.885 24.348 5.05 10.487 7.12 23.674 6.846 38.674-.5 27.368-8.862 60.148-17.2 91.362l-36.864-9.88-51.232 153.712-42.69.11-1.23 18.69 57.402-.146 49.914-149.758 37.946 10.166 2.42-9.025c9.022-33.677 19.603-71.135 20.22-104.89.31-16.876-1.89-32.994-8.693-47.124-5.016-10.417-12.696-19.57-23.065-26.622 10.814-11.607 19.228-27.125 23.637-43.58 11.288-42.13-6.228-85.52-42.38-95.21l-.003-.003c-1.106-.296-3.297-1.274-6.81-1.744h-.008l-2.838-.38-.295.146c-1.09-.082-2.185-.226-3.27-.244zm-349.32.46c-4.49.056-9.02.665-13.538 1.876-.095.026-.327.068-.44.094l-.575-.574-5.76 2.377h-.002C27.32 36.99 13.11 77.635 23.69 117.12c4.574 17.073 13.46 32.977 24.845 44.67-9.328 6.978-16.34 15.908-21.053 25.99-6.507 13.924-8.973 29.83-9.11 46.6-.27 33.543 8.753 71.01 17.82 104.845l2.42 9.027 40.02-10.727 51.11 149.454 60.46.153-1.39-18.694-45.7-.116-52.446-153.37-38.73 10.378c-8.028-30.892-15.098-63.467-14.875-90.8.122-14.997 2.417-28.276 7.354-38.84 4.937-10.56 12.24-18.566 23.865-24.15l14.298-6.87-12.94-9.176c-11.456-8.122-23.12-25.39-27.896-43.215-8.66-32.315 3.867-62.596 24.653-71.188l.025-.01c.244-.1 1.86-.42 4.486-1.12h.002l.002-.003c2.966-.796 6.005-1.18 9.072-1.175 21.47.027 44.263 19.06 52.344 49.223 4.66 17.392 3.46 37.92-2.035 50.517l-6.436 14.76 16.01-1.734c13.355-1.447 23.684 1.234 32.868 7.016 4.285-4.866 9.108-9.17 14.46-12.742-.73-.536-1.464-1.062-2.212-1.572-9.55-6.512-20.777-10.598-33.283-11.522 3.562-15.46 3.09-33.105-1.318-49.56-9.878-36.864-39.338-63.538-70.77-63.14z", 512, 512)
	addFileType(".adq", fileSvg, 384, 512) // TODO: Create icon
	addFileType(".adm", fileSvg, 384, 512) // TODO: Create icon
	addFileType(".eqp", fileSvg, 384, 512) // TODO: Create icon
	addFileType(".eqm", fileSvg, 384, 512) // TODO: Create icon
	addFileType(".skl", fileSvg, 384, 512) // TODO: Create icon
	addFileType(".spl", "M103.432 17.844c-1.118.005-2.234.032-3.348.08-2.547.11-5.083.334-7.604.678-20.167 2.747-39.158 13.667-52.324 33.67-24.613 37.4 2.194 98.025 56.625 98.025.536 0 1.058-.012 1.583-.022v.704h60.565c-10.758 31.994-30.298 66.596-52.448 101.43-2.162 3.4-4.254 6.878-6.29 10.406l34.878 35.733-56.263 9.423c-32.728 85.966-27.42 182.074 48.277 182.074v-.002l9.31.066c23.83-.57 46.732-4.298 61.325-12.887 4.174-2.458 7.63-5.237 10.467-8.42h-32.446c-20.33 5.95-40.8-6.94-47.396-25.922-8.956-25.77 7.52-52.36 31.867-60.452 5.803-1.93 11.723-2.834 17.565-2.834v-.406h178.33c-.57-44.403 16.35-90.125 49.184-126 23.955-26.176 42.03-60.624 51.3-94.846l-41.225-24.932 38.272-6.906-43.37-25.807h-.005l.002-.002.002.002 52.127-8.85c-5.232-39.134-28.84-68.113-77.37-68.113C341.14 32.26 222.11 35.29 149.34 28.496c-14.888-6.763-30.547-10.723-45.908-10.652zm.464 18.703c13.137.043 27.407 3.804 41.247 10.63l.033-.07c4.667 4.735 8.542 9.737 11.68 14.985H82.92l10.574 14.78c10.608 14.83 19.803 31.99 21.09 42.024.643 5.017-.11 7.167-1.814 8.836-1.705 1.67-6.228 3.875-15.99 3.875-40.587 0-56.878-44.952-41.012-69.06C66.238 46.64 79.582 39.22 95.002 37.12c2.89-.395 5.863-.583 8.894-.573zM118.5 80.78h46.28c4.275 15.734 3.656 33.07-.544 51.51H131.52c1.9-5.027 2.268-10.574 1.6-15.77-1.527-11.913-7.405-24.065-14.62-35.74zm101.553 317.095c6.44 6.84 11.192 15.31 13.37 24.914 3.797 16.736 3.092 31.208-1.767 43.204-4.526 11.175-12.576 19.79-22.29 26h237.19c14.448 0 24.887-5.678 32.2-14.318 7.312-8.64 11.2-20.514 10.705-32.352-.186-4.473-.978-8.913-2.407-13.18l-69.91-8.205 42.017-20.528c-8.32-3.442-18.64-5.537-31.375-5.537H220.053zm-42.668.506c-1.152-.003-2.306.048-3.457.153-2.633.242-5.256.775-7.824 1.63-15.11 5.02-25.338 21.54-20.11 36.583 3.673 10.57 15.347 17.71 25.654 13.938l1.555-.57h43.354c.946-6.36.754-13.882-1.358-23.192-3.71-16.358-20.543-28.483-37.815-28.54z", 512, 512)
	addFileType(".not", "M224 136V0H24C10.7 0 0 10.7 0 24v464c0 13.3 10.7 24 24 24h336c13.3 0 24-10.7 24-24V160H248c-13.2 0-24-10.8-24-24zm64 236c0 6.6-5.4 12-12 12H108c-6.6 0-12-5.4-12-12v-8c0-6.6 5.4-12 12-12h168c6.6 0 12 5.4 12 12v8zm0-64c0 6.6-5.4 12-12 12H108c-6.6 0-12-5.4-12-12v-8c0-6.6 5.4-12 12-12h168c6.6 0 12 5.4 12 12v8zm0-72v8c0 6.6-5.4 12-12 12H108c-6.6 0-12-5.4-12-12v-8c0-6.6 5.4-12 12-12h168c6.6 0 12 5.4 12 12zm96-114.1v6.1H256V0h6.1c6.4 0 12.5 2.5 17 7l97.9 98c4.5 4.5 7 10.6 7 16.9z", 384, 512)
	addFileType(".pdf", "M181.9 256.1c-5-16-4.9-46.9-2-46.9 8.4 0 7.6 36.9 2 46.9zm-1.7 47.2c-7.7 20.2-17.3 43.3-28.4 62.7 18.3-7 39-17.2 62.9-21.9-12.7-9.6-24.9-23.4-34.5-40.8zM86.1 428.1c0 .8 13.2-5.4 34.9-40.2-6.7 6.3-29.1 24.5-34.9 40.2zM248 160h136v328c0 13.3-10.7 24-24 24H24c-13.3 0-24-10.7-24-24V24C0 10.7 10.7 0 24 0h200v136c0 13.2 10.8 24 24 24zm-8 171.8c-20-12.2-33.3-29-42.7-53.8 4.5-18.5 11.6-46.6 6.2-64.2-4.7-29.4-42.4-26.5-47.8-6.8-5 18.3-.4 44.1 8.1 77-11.6 27.6-28.7 64.6-40.8 85.8-.1 0-.1.1-.2.1-27.1 13.9-73.6 44.5-54.5 68 5.6 6.9 16 10 21.5 10 17.9 0 35.7-18 61.1-61.8 25.8-8.5 54.1-19.1 79-23.2 21.7 11.8 47.1 19.5 64 19.5 29.2 0 31.2-32 19.7-43.4-13.9-13.6-54.3-9.7-73.6-7.2zM377 105L279 7c-4.5-4.5-10.6-7-17-7h-6v128h128v-6.1c0-6.3-2.5-12.4-7-16.9zm-74.1 255.3c4.1-2.7-2.5-11.9-42.8-9 37.1 15.8 42.8 9 42.8 9z", 384, 512)
	imageFileSvg := "M384 121.941V128H256V0h6.059a24 24 0 0 1 16.97 7.029l97.941 97.941a24.002 24.002 0 0 1 7.03 16.971zM248 160c-13.2 0-24-10.8-24-24V0H24C10.745 0 0 10.745 0 24v464c0 13.255 10.745 24 24 24h336c13.255 0 24-10.745 24-24V160H248zm-135.455 16c26.51 0 48 21.49 48 48s-21.49 48-48 48-48-21.49-48-48 21.491-48 48-48zm208 240h-256l.485-48.485L104.545 328c4.686-4.686 11.799-4.201 16.485.485L160.545 368 264.06 264.485c4.686-4.686 12.284-4.686 16.971 0L320.545 304v112z"
	addFileType(".png", imageFileSvg, 384, 512)
	addFileType(".jpg", imageFileSvg, 384, 512)
	addFileType(".jpeg", imageFileSvg, 384, 512)
	addFileType(".webp", imageFileSvg, 384, 512)
	addFileType(".gif", imageFileSvg, 384, 512)
}

func addFileType(ext, svgPath string, width, height float32) {
	p, err := unison.NewSvgPath(geom32.NewSize(width, height), svgPath)
	jot.FatalIfErr(err)
	fileTypes[ext] = p
}
