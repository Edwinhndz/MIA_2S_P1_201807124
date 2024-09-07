package filesystem

import (
	"fmt"
	"strconv"
	"strings"
)

// DividirComando recibe un comando y lo divide en un arreglo de strings
func DividirComando(comando string) string {
	var respuesta string
	//se divide el comando en un arreglo de strings por el enter
	comandos := strings.Split(comando, "\n")
	//se recorre el arreglo de strings
	for i := 0; i < len(comandos); i++ {
		//imprime el comando
		fmt.Println("Comando: " + comandos[i])
		//se analiza el comando
		respuesta += AnalizarComando(comandos[i])
	}
	return respuesta
}

// AnalizarComando recibe un comando y lo analiza
func AnalizarComando(comando string) string {
	var respuesta string
	//se divide el comando en un arreglo de strings por el espacio
	comandoSeparado := strings.Split(comando, " ")
	//Si encuentra el # en la posicion 0, es un comentario
	if strings.Contains(comandoSeparado[0], "#") {
		//imprime el comentario
		fmt.Println("Comentario: ")
		//Eliminiar el #
		comandoSeparado[0] = strings.Replace(comandoSeparado[0], "#", "", -1)
		respuesta += "Comentario: "
		//se recorre el arreglo de strings
		for i := 0; i < len(comandoSeparado); i++ {
			respuesta += comandoSeparado[i] + " "
		}
		respuesta += "\n"
		fmt.Println(respuesta)
	} else {
		//Si no es un comentario, entonces es un comando
		//Iterar sobre el comando
		for _, valor := range comandoSeparado {
			//el primer valor del comando lo pasamos a minusculas
			valor = strings.ToLower(valor)
			//Si el valor es mkdisk, entonces es un comando para crear un disco
			if valor == "mkdisk" {
				fmt.Println("------------------|Comando mkdisk|------------------")
				respuesta += "------------------|Comando mkdisk|------------------\n"
				//Analizar el comando mkdisk
				respuesta += AnalizarMkdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "rmdisk" {
				fmt.Println("------------------|Comando rmdisk|------------------")
				respuesta += "Ejecutando rmdisk\n"
				//Analizar el comando rmdisk
				//respuesta += AnalizarRmdisk(comandoSeparado)
				//Pasar a string el comando separado
				//comandoSeparadoString := strings.Join(comandoSeparado, " ")
				//respuesta += AnalizarComando(comandoSeparadoString)
				//return respuesta
			} else if valor == "fdisk" {
				fmt.Println("------------------|Comando fdisk|------------------")
				respuesta += "------------------|Comando fdisk------------------\n"
				//Analizar el comando fdisk
				respuesta += AnalizarFdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "mount" {
				fmt.Println("------------------|Comando mount|------------------")
				respuesta += "------------------|Comando mount|------------------\nParametros:\n"
				//Analizar Comando Mount
				respuesta += analizarMount(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "mkfs" {
				fmt.Println("------------------|Comando mkfs|------------------")
				respuesta += "------------------|Comando mkfs|------------------\nParametros:\n"
				//Analizar Comando Mkfs
				respuesta += analizarMkfs(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "login" {
				fmt.Println("------------------|Comando login|------------------")
				respuesta += "------------------|Comando login|------------------\nParametros:\n"	
				//Analizar Comando Login
				respuesta += analizarLogin(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta

			} else if valor == "rep" {
				fmt.Println("------------------|Comando rep|------------------")
				respuesta += "------------------|Comando rep|------------------\n\n"
				//Analizar Comando Rep
				respuesta += analizarRep(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta

			} else if valor == "\n" {
				continue
			} else if valor == "\r" {
				continue
			} else if valor == "\t" {
				continue
			} else if valor == "" {
				continue
			} else {
				//Si no es ningun comando, entonces es un error
				fmt.Println("Error: Comando no reconocido")
				respuesta += "Error: Comando no reconocido\n"
			}
		}
	}
	return respuesta
}

// AnalizarMkdisk recibe un comando mkdisk y lo analiza
func AnalizarMkdisk(comando *[]string) string {
	//mkdisk -unit=M -path="/home 1/mis discos/Disco3.mia"
	//0 		1     2     3"/home/mis     4         5discos/Disco3.mia"
	var respuesta string
	//Eliminar el primer valor del comando
	*comando = (*comando)[1:]
	//-size=5 -unit=M -path="/home/mis discos/Disco3.mia"
	//Booleanos para saber si se encontro el size, unit, fit, path
	var size, unit, path, fit bool
	//Variables para guardar el valor del size, unit, fit, path
	var valorSize, valorUnit, valorFit, valorPath string
	//Iterar sobre el comando
	valorFit = "f"
	valorUnit = "m"
	for _, valor := range *comando {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-size" {
			size = true
			valorSize = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-unit" {
			unit = true
			valorUnit = banderaValor
			valorUnit = strings.ToLower(valorUnit)
			*comando = (*comando)[1:]
		} else if bandera == "-fit" {
			fit = true
			valorFit = banderaValor
			valorFit = strings.ToLower(valorFit)
			*comando = (*comando)[1:]
		} else if bandera == "-path" {
			path = true
			//Verificar si el path tiene comillas
			//-path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
				//Eliminar las comillas del inicio
				banderaValor = strings.Replace(banderaValor, "\"", "", -1)
				//Eliminar el primer valor del comandoSeparado
				*comando = (*comando)[1:]
				//Iterar sobre el comando
				Contador := 0
				for _, valor := range *comando {
					//Si el valor contiene comillas
					if strings.Contains(valor, "\"") {
						//Eliminar las comillas del final
						valor = strings.Replace(valor, "\"", "", -1)
						//Agregar el valor al path
						valorPath += valor
						break
					} else {
						//Agregar el valor al path
						valorPath += valor + " "
						Contador++
					}
				}
				//Eliminar los valores del comando
				*comando = (*comando)[Contador:]
			} else {
				valorPath = banderaValor
				*comando = (*comando)[1:]
			}
		} else {
			fmt.Println("Error: Parametro no reconocida")
			respuesta += "\nError: Parametro no reconocida\n"
		}

	}
	if !size {
		fmt.Println("Error: Falta el parametro size")
		respuesta += "\nError: Falta el parametro size\n"
		return respuesta
	} else if !path {
		fmt.Println("Error: Falta el parametro path")
		respuesta += "\nError: Falta el parametro path\n"
		return respuesta
	} else {
		if fit {
			if valorFit != "bf" && valorFit != "ff" && valorFit != "wf" {
				fmt.Println("Error: Fit no reconocido")
				respuesta += "\nError: Fit no reconocido\n"
				return respuesta
			} else {
				if valorFit == "bf" {
					valorFit = "b"
				} else if valorFit == "ff" {
					valorFit = "f"
				} else if valorFit == "wf" {
					valorFit = "w"
				}
			}
		}
		if unit {
			if valorUnit != "k" && valorUnit != "m" {
				fmt.Println("Error: Unit no reconocido")
				respuesta += "\nError: Unit no reconocido\n"
				return respuesta
			}
		}
		//Pasar a int el size
		sizeInt, err := strconv.Atoi(valorSize)
		if err != nil {
			fmt.Println("Error: Size no es un numero")
			respuesta += "\nError: Size no es un numero\n"
			return respuesta
		}
		//Verificar que el size sea mayor a 0
		if sizeInt <= 0 {
			fmt.Println("Error: Size debe ser mayor a 0")
			respuesta += "\nError: Size debe ser mayor a 0\n"
			return respuesta
		}
		//Imprimir los valores
		fmt.Println("Size: " + valorSize)
		fmt.Println("Unit: " + valorUnit)
		fmt.Println("Fit: " + valorFit)
		fmt.Println("Path: " + valorPath)
		//imprime los valores en el frontend
		respuesta += "Parametros:\n"
		respuesta += "Size: " + valorSize + "\n"
		respuesta += "Unit: " + valorUnit + "\n"
		respuesta += "Fit: " + valorFit + "\n"
		respuesta += "Path: " + valorPath + "\n\n"

		//Llamar a la funcion para crear el disco
		respuesta += CrearDiscos(sizeInt, valorUnit, valorFit, valorPath)
		return respuesta
	}
	
}

// AnalizarFdisk recibe un comando fdisk y lo analiza
func AnalizarFdisk(comando *[]string) string {
	//fdisk -Size=300 -path=/home/Disco1.mia -name=Particion1
	*comando = (*comando)[1:]
	//respuesta
	var respuesta string
	//Booleanos para saber si se encontro el size, unit, fit, path
	var size, unit, path, name, typePart, fit bool
	//Variables para guardar el valor del size, unit, fit, path
	var valorSize, valorUnit, valorFit, valorPath, valorName, valorTypePart string
	valorFit = "f"
	valorUnit = "k"
	valorTypePart = "p"
	//Iterar sobre el comando
	for _, valor := range *comando {
		//Obtener la bandera
		bandera := obtenerBandera(valor)
		//Obtener el valor de la bandera
		banderaValor := obtenerValor(valor)
		//Si la bandera es -size
		if bandera == "-size" {
			size = true
			valorSize = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-unit" {
			unit = true
			valorUnit = banderaValor
			valorUnit = strings.ToLower(valorUnit)
			*comando = (*comando)[1:]
		} else if bandera == "-fit" {
			fit = true
			valorFit = banderaValor
			valorFit = strings.ToLower(valorFit)
			*comando = (*comando)[1:]
		} else if bandera == "-name" {
			name = true
			valorName = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-type" {
			typePart = true
			valorTypePart = banderaValor
			valorTypePart = strings.ToLower(valorTypePart)
			*comando = (*comando)[1:]
		} else if bandera == "-path" {
			path = true
			//Verificar si el path tiene comillas
			//-path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
				//Eliminar las comillas del inicio
				banderaValor = strings.Replace(banderaValor, "\"", "", -1)
				//Eliminar el primer valor del comandoSeparado
				*comando = (*comando)[1:]
				//Iterar sobre el comando
				Contador := 0
				for _, valor := range *comando {
					//Si el valor contiene comillas
					if strings.Contains(valor, "\"") {
						//Eliminar las comillas del final
						valor = strings.Replace(valor, "\"", "", -1)
						//Agregar el valor al path
						valorPath += valor
						break
					} else {
						//Agregar el valor al path
						valorPath += valor + " "
						Contador++
					}
				}
				//Eliminar los valores del comando
				*comando = (*comando)[Contador:]
			} else {
				valorPath = banderaValor
				*comando = (*comando)[1:]
			}
		} else {
			fmt.Println("Error: Parametro no reconocida")
			respuesta += "\nError: Parametro no reconocida\n"
		}
	}
	//Obligatorios: name, path y size
	if !name {
		fmt.Println("Error: Falta el parametro name")
		respuesta += "\nError: Falta el parametro name\n"
		return respuesta
	} else if !path {
		fmt.Println("Error: Falta el parametro path")
		respuesta += "\nError: Falta el parametro path\n"
		return respuesta
	} else if !size {
		fmt.Println("Error: Falta el parametro size")
		respuesta += "\nError: Falta el parametro size\n"
		return respuesta
	} else {
		//Opcionales: unit, fit, type
		if fit {
			if valorFit != "bf" && valorFit != "ff" && valorFit != "wf" {
				fmt.Println("Error: Fit no reconocido")
				respuesta += "\nError: Fit no reconocido\n"
				return respuesta
			} else {
				if valorFit == "bf" {
					valorFit = "b"
				} else if valorFit == "ff" {
					valorFit = "f"
				} else if valorFit == "wf" {
					valorFit = "w"
				}
			}
		}
		if unit {
			if valorUnit != "k" && valorUnit != "m" && valorUnit != "b" {
				fmt.Println("Error: Unit no reconocido")
				respuesta += "\nError: Unit no reconocido\n"
				return respuesta
			}
		}
		if typePart {
			if valorTypePart != "p" && valorTypePart != "e" && valorTypePart != "l" {
				fmt.Println("Error: Type no reconocido")
				respuesta += "\nError: Type no reconocido\n"
				return respuesta
			}
		}
		//Pasar a int el size
		sizeInt, err := strconv.Atoi(valorSize)
		if err != nil {
			fmt.Println("Error: Size no es un numero")
			respuesta += "\nError: Size no es un numero\n"
			return respuesta
		}
		//Verificar que el size sea mayor a 0
		if sizeInt <= 0 {
			fmt.Println("Error: Size debe ser mayor a 0")
			respuesta += "\nError: Size debe ser mayor a 0\n"
			return respuesta
		}
		//Imprimir los valores
		fmt.Println("Size: " + valorSize)
		fmt.Println("Unit: " + valorUnit)
		fmt.Println("Fit: " + valorFit)
		fmt.Println("Path: " + valorPath)
		fmt.Println("Name: " + valorName)
		fmt.Println("Type: " + valorTypePart)

		respuesta += "Parametros:\n"
		respuesta += "Size: " + valorSize + "\n"
		respuesta += "Unit: " + valorUnit + "\n"
		respuesta += "Fit: " + valorFit + "\n"
		respuesta += "Path: " + valorPath + "\n"
		respuesta += "Name: " + valorName + "\n"
		respuesta += "Type: " + valorTypePart + "\n\n"
	
		//Llamar a la funcion para crear la particion
		respuesta += Fdisk(sizeInt, valorUnit, valorFit, valorPath, valorName, valorTypePart)
		return respuesta
	}
}

func analizarMount(comandoSeparado *[]string) string {
	//respuesta
	var respuesta string
	//mount -driveletter=A -name=Part1 #id=A118
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaPath, banderaName bool
	//Variables para almacenar los valores de los parametros
	var pathValor, nameValor string
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-path" {
			banderaPath = true
			//Verificar si el path tiene comillas
			//-path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
				//Eliminar las comillas del inicio
				banderaValor = strings.Replace(banderaValor, "\"", "", -1)
				//Eliminar el primer valor del comandoSeparado
				*comandoSeparado = (*comandoSeparado)[1:]
				//Iterar sobre el comando
				Contador := 0
				for _, valor := range *comandoSeparado {
					//Si el valor contiene comillas
					if strings.Contains(valor, "\"") {
						//Eliminar las comillas del final
						valor = strings.Replace(valor, "\"", "", -1)
						//Agregar el valor al path
						pathValor += valor
						break
					} else {
						//Agregar el valor al path
						pathValor += valor + " "
						Contador++
					}
				}
				//Eliminar los valores del comando
				*comandoSeparado = (*comandoSeparado)[Contador:]
			} else {
				pathValor = banderaValor
				*comandoSeparado = (*comandoSeparado)[1:]
			}
		} else if bandera == "-name" {
			banderaName = true
			nameValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("\nParametro no reconocido: ", bandera)
			respuesta += "\nParametro no reconocido: " + bandera + "\n"
		}
	}
	//Obligatorios: -driveletter, -name
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaPath {
		fmt.Println("\nEl parametro -path es obligatorio")
		respuesta += "\nEl parametro -path es obligatorio\n"
		return respuesta
	} else if !banderaName {
		fmt.Println("\nEl parametro -name es obligatorio")
		respuesta += "\nEl parametro -name es obligatorio\n"
		return respuesta
	} else {
		//Imprimir los valores de los parametros
		fmt.Println("\nPath: ", pathValor)
		respuesta += "Path: " + pathValor + "\n"
		fmt.Println("Name: ", nameValor)
		respuesta += "Name: " + nameValor + "\n"
		//Llamar a la funcion para montar la particion
		respuesta += MountPartition(pathValor, nameValor)
		return respuesta
	}
}

func analizarMkfs(comandoSeparado *[]string) string {
	// mkfs -type=full -id=B145 -fs=3fs
	//respuesta
	var respuesta string
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaType, banderaId, banderaFs bool
	//Variables para almacenar los valores de los parametros
	var typeValor, idValor, fsValor string
	typeValor = "full"
	fsValor = "2FS"
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-type" {
			banderaType = true
			typeValor = banderaValor
			typeValor = strings.ToLower(typeValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-id" {
			banderaId = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-fs" {
			banderaFs = true
			fsValor = banderaValor
			fsValor = strings.ToLower(fsValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("\nParametro no reconocido: ", bandera)
			respuesta += "\nParametro no reconocido: " + bandera + "\n"
		}
	}
	//Obligatorios: -id
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaId {
		fmt.Println("El parametro -id es obligatorio")
		respuesta += "El parametro -id es obligatorio\n"
		return respuesta
	} else {
		//Verificar si se ingresaron los parametros aceptados
		if banderaType {
			if typeValor != "full" {
				fmt.Println("El valor del parametro -type no es valido")
				respuesta += "El valor del parametro -type no es valido\n"
				return respuesta
			}
		}
		if banderaFs {
			if fsValor != "2FS" && fsValor != "3fs" {
				fmt.Println("El valor del parametro -fs no es valido")
				respuesta += "El valor del parametro -fs no es valido\n"
				return respuesta
			}
		}
		//Imprimir los valores de los parametros
		fmt.Println("Type: ", typeValor)
		fmt.Println("Id: ", idValor)
		fmt.Println("Fs: ", fsValor)
		//imprime los valores en el frontend
		respuesta += "Type: " + typeValor + "\n"
		respuesta += "Id: " + idValor + "\n"
		respuesta += "Fs: " + fsValor + "\n\n"

		//Llamar a la funcion para formatear la particion
		respuesta += Mkfs(typeValor, idValor, fsValor)
		return respuesta
	}
}

func analizarLogin(comandoSeparado *[]string) string {
	//mount -driveletter=A -name=Part1 #id=A118
	//respuesta
	var respuesta string
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaUser, banderaPass, banderaId bool
	//Variables para almacenar los valores de los parametros
	var userValor, passValor, idValor string
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-user" {
			banderaUser = true
			userValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-pass" {
			banderaPass = true
			passValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-id" {
			banderaId = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
			respuesta += "\nParametro no reconocido: " + bandera + "\n"
		}
	}
	//Obligatorios: -user, -pass, -id
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaUser {
		fmt.Println("El parametro -user es obligatorio")
		respuesta += "\nEl parametro -user es obligatorio\n"	
	}
	if !banderaPass {
		fmt.Println("El parametro -pass es obligatorio")
		respuesta += "\nEl parametro -pass es obligatorio\n"
	}	
	if !banderaId {	
		fmt.Println("El parametro -id es obligatorio")
		respuesta += "\nEl parametro -id es obligatorio\n"
	} else {
		//Imprimir los valores de los parametros
		fmt.Println("User: ", userValor)
		fmt.Println("Pass: ", passValor)
		fmt.Println("Id: ", idValor)
		respuesta += "User: "+userValor+"\n"
		respuesta += "Pass: "+passValor+"\n"
		respuesta += "Id: "+idValor+ "\n"
		//Llamar a la funcion para montar la particion

		respuesta += Login(userValor, passValor, idValor)
	}
	return respuesta
}

func analizarRep(comandoSeparado *[]string) string {
	//rep -id=A118 -path=/home/user/reports/report2.jpg -name=disk	
	//respuesta
	var respuesta string
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaId, banderaPath, banderaName, banderaRuta bool
	//Variables para almacenar los valores de los parametros
	var idValor, pathValor, nameValor, rutaValor string
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-id" {
			banderaId = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-path" {	
			banderaPath = true
			pathValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-name" {
			banderaName = true
			nameValor = banderaValor
			nameValor = strings.ToLower(nameValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if  bandera == "-ruta" {
			banderaRuta = true
			rutaValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("\nParametro no reconocido: ", bandera)
			respuesta += "\nParametro no reconocido: " + bandera + "\n"
		}
	}
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaId {
		fmt.Println("El parametro -id es obligatorio")
		respuesta += "El parametro -id es obligatorio\n"
	}
	if !banderaPath {
		fmt.Println("El parametro -path es obligatorio")
		respuesta += "El parametro -path es obligatorio\n"
	}
	if !banderaName {
		fmt.Println("El parametro -name es obligatorio")
		respuesta += "El parametro -name es obligatorio\n"
	}

	// Verificar si -ruta es obligatorio en función del valor de -name
	if (nameValor == "file" || nameValor == "ls") && !banderaRuta {
		fmt.Println("El parametro -ruta es obligatorio cuando -name es 'file' o 'ls'")
		respuesta += "El parametro -ruta es obligatorio cuando -name es 'file' o 'ls'\n"
	}
	
	// Procesar comando si todo es correcto
	if len(respuesta) == 0 { // Si no hay errores
		fmt.Println("Id: ", idValor)
		fmt.Println("Path: ", pathValor)
		fmt.Println("Name: ", nameValor)
		fmt.Println("Ruta: ", rutaValor)
		//imprime los valores en el frontend
		respuesta += "Parametros:\n"
		respuesta += "Id: " + idValor + "\n"
		respuesta += "Path: " + pathValor + "\n"
		respuesta += "Name: " + nameValor + "\n"
		respuesta += "Ruta: " + rutaValor + "\n\n"
		if nameValor == "tree" {
			RepTree(idValor, pathValor)
		} else if nameValor == "disk" {
			ReporteDisk(idValor, pathValor)
		} else if nameValor == "sb" {
			ReporteSB(idValor, pathValor)
		}
	}
	return respuesta
}


func obtenerBandera(bandera string) string {
	//-size
	var banderaValor string
	for _, valor := range bandera {
		if valor == '=' {
			break
		}
		banderaValor += string(valor)
	}
	banderaValor = strings.ToLower(banderaValor)
	return banderaValor
}

func obtenerValor(bandera string) string {
	var banderaValor string
	var boolBandera bool
	for _, valor := range bandera {
		if boolBandera {
			banderaValor += string(valor)
		}
		if valor == '=' {
			boolBandera = true
		}
	}
	return banderaValor
}

