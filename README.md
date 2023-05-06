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

| ![image](https://user-images.githubusercontent.com/94043894/236626257-7fb68cf0-89e7-4230-885f-f6f62b95490b.gif) |
| :-------------------------------------------------------------------------------------------------------------: |
|                                          gif, not much use though. 💩                                           |

### 👇 Install

1. `releases` download.
2. build yourself, makefile in [cli/makefile](cli/makefile).
3. maybe used as sdk ? `go get github.com/orzation/bobibo`.
4. `AUR` use `yay/paru -S bobibo`.

### 🍰 How2use

`bobibo [-option] /path/to/image.png `

options:

- `-v` enable reverse the character color.
- `-g` enable gif mode(test), print every frame of gif image.
- `-s value` set the scale for images(value default 0.5, (0, +)).
- `-t value` set the threshold of binarization(value default generate by OTSU, [-1, 255]).

> use `bobibo -h` to print options.
> use `bobibo version` to print version.

### ⚙️ Contribute

> hope so 💩

1. fork
2. do your things
3. pull request

### 📄 License

GPLV3.0
