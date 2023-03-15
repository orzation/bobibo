## 🐱 Bobibo

### 😗 Introduce

Bobibo is a simple cli-tool, it can convert pictures(jpeg, png, even gif) to ascii arts using
braille unicode.

> I wrote it, cause i need some anime ascii arts. 💩

### 🎞️ Screenshot

| ![image](https://user-images.githubusercontent.com/94043894/223673376-f67f030c-305f-4dd1-beee-301a8da79b5d.png) |
| :-------------------------------------------------------------------------------------------------------------: |
|                                                bobibobibobobibo                                                 |

| ![image](https://user-images.githubusercontent.com/94043894/223674513-ed33023d-9181-4fe6-bf7e-cd059bfd0ba3.png) |
| :-------------------------------------------------------------------------------------------------------------: |
| ![image](https://user-images.githubusercontent.com/94043894/223675190-ecbd20a6-cf49-40a0-a36d-d7bf6b0a75ff.png) |
|                                   reverse when your background is too light.                                    |

| ![image](https://user-images.githubusercontent.com/94043894/223677662-d27bc50a-3364-461f-bad4-ba7f0c4b8df9.png) |
| :-------------------------------------------------------------------------------------------------------------: |
|                                          gif, not much use though. 💩                                           |

### 👇 Install

1. `releases` download.
2. build yourself, makefile in [cli/makefile](cli/makefile).
3. maybe used as sdk ? `go get github.com/orzation/bobibo`.
4. `AUR` use `yay/paru -S bobibo`.

### 🍰 How2use

`bobibo /path/to/image.png [-option]`

options:

- `-r` enable reverse the character color.
- `-g` enable gif mode, print every frame of gif image.
- `-s value` set the scale for images(value default 0.5, (0, +)).
- `-t value` set the threshold of binarization(value default generate by OTSU, [0, 255]).

> use `bobibo help` to print options.
> use `bobibo version` to print version.

### ⚙️ Contribute

> hope so 💩

1. fork
2. do your things
3. pull request

### 📄 License

GPLV3.0
