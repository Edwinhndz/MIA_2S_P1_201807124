package filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func ReporteDisk(idValor string, pathValor string) {
	//Abrir el disco A

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	//Buscar la particion montada con el ID
	Particion := VerificarParticionMontada(idValor)
	if Particion == -1 {
		fmt.Println("No se encontro la particion montada con el ID: ", idValor)
		return
	} else {
		fmt.Println("Se encontro la particion montada con el ID: ", idValor)
	}
	MountActual := particionesMontadas[Particion]
	//Abrir el disco
	archivo, err := os.OpenFile(""+MountActual.Path+"", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	disk := NewMBR()
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		return
	}
	sizeMBR := int(disk.Mbr_tamano)
	libre := int(disk.Mbr_tamano)

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte Disk \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "node[shape=record, color=lightgrey]a0[label=\"MBR"

	if disk.Mbr_partition_1.Part_size != 0 {
		libre -= int(disk.Mbr_partition_1.Part_size)
		Dot += "|"
		if disk.Mbr_partition_1.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_1.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_1.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_1.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition_2.Part_size != 0 {
		libre -= int(disk.Mbr_partition_2.Part_size)
		Dot += "|"
		if disk.Mbr_partition_2.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_2.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_2.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_2.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition_3.Part_size != 0 {
		libre -= int(disk.Mbr_partition_3.Part_size)
		Dot += "|"
		if disk.Mbr_partition_3.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_3.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_3.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_3.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition_4.Part_size != 0 {
		libre -= int(disk.Mbr_partition_4.Part_size)
		Dot += "|"
		if disk.Mbr_partition_4.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_4.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_4.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			eb := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_4.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &eb)
			if err != nil {
				fmt.Println("Error al leer el EBR del disco: ", err)
				return
			}
			if eb.Part_size != 0 {
				Dot += "|{"
				PrimerEbr := true
				for {
					if !PrimerEbr {
						Dot += "|EBR"
					} else {
						PrimerEbr = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					porcentaje := (float64(eb.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libre -= int(eb.Part_size)

					Desplazamiento += int(eb.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					binary.Read(archivo, binary.LittleEndian, &eb)
					if eb.Part_size == 0 {
						break
					}
				}
				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if libre > 0 {
		Dot += "|Libre"
		porcentaje := (float64(libre) * float64(100)) / float64(sizeMBR)
		Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
	}
	Dot += "\"];\n}"

	//Quitar la extension al archivo (pdf, etc, )

	//Crear el archivo .dot
	DotName := "Reportes/ReporteDisk.dot"
	archivoDot, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		return
	}
	defer archivoDot.Close()
	_, err = archivoDot.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		return
	}
	//Generar la imagen
	cmd := exec.Command("dot", "-T", "png", DotName, "-o", "Reportes/ReporteDisk.png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte generado con exito")
	//se escribe en el front

}

func RepTree(idValor string, pathValor string) {
	//Crear el directorio si no existe
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	//Buscar la particion montada con el ID
	Particion := VerificarParticionMontada(idValor)
	if Particion == -1 {
		fmt.Println("No se encontro la particion montada con el ID: ", idValor)
		return
	}
	MountActual := particionesMontadas[Particion]
	//Abrir el disco
	archivo, err := os.OpenFile(""+MountActual.Path+"", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	archivo.Seek(int64(MountActual.Start), 0)
	//Leer el superbloque
	sb := NewSuperBlock()
	err = binary.Read(archivo, binary.LittleEndian, &sb)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	//Buscar el inodo raiz
	raiz := NewInodes()
	archivo.Seek(int64(sb.S_inode_start), 0)
	binary.Read(archivo, binary.LittleEndian, &raiz)
	Dot := "digraph H {\n"
	Dot += "node [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"1\"];\n"
	Dot += "node [shape=plaintext];\n"
	Dot += "graph [bb=\"0,0,352,154\"];\n"
	Dot += "rankdir=LR;\n"
	Dot += RecursivoTree(raiz, sb, archivo, 0)
	Dot += "}"

	//Quitar la extension al archivo (pdf, etc, )
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
	archivoDot, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		return
	}
	defer archivoDot.Close()
	_, err = archivoDot.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		return
	}

	//Generar la imagen
	//Quitar el punto
	extensionSinPunto := strings.TrimPrefix(extension, ".")
	//Correr con todos los permisos
	cmd := exec.Command("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	fmt.Println("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		return
	}

	fmt.Println("Reporte tree generado con exito")

}

func RecursivoTree(inodo Inodes, sb SuperBlock, archivo *os.File, numeroInodo int) string {
	Dot := "Inodo" + strconv.Itoa(numeroInodo) + "[label = <\n"
	Dot += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	Dot += "<tr><td bgcolor=\"lightgrey\">Inodo" + strconv.Itoa(numeroInodo) + "</td></tr>\n"
	Dot += "<tr><td>i_uid</td><td>" + strconv.Itoa(int(inodo.I_uid)) + "</td></tr>\n"
	Dot += "<tr><td>i_gid</td><td>" + strconv.Itoa(int(inodo.I_gid)) + "</td></tr>\n"
	Dot += "<tr><td>i_size</td><td>" + strconv.Itoa(int(inodo.I_size)) + "</td></tr>\n"
	Dot += "<tr><td>i_atime</td><td>" + string(inodo.I_atime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_ctime</td><td>" + string(inodo.I_ctime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_mtime</td><td>" + string(inodo.I_mtime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_type</td><td>" + string(inodo.I_type[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_perm</td><td>" + strconv.Itoa(int(inodo.I_perm)) + "</td></tr>\n"
	//Recorrer los bloques
	Contador := 0
	for _, i := range inodo.I_block {
		Dot += "<tr><td>i_block" + strconv.Itoa(Contador+1) + "</td><td port='" + strconv.Itoa(Contador+1) + "'>" + strconv.Itoa(int(i)) + "</td></tr>\n"
		Contador++
	}
	Dot += "</table>>];\n"
	//Recorrer los bloques
	Contador = 0
	for _, i := range inodo.I_block {
		if i != -1 {
			//Leer el bloque
			Dot += "Inodo" + strconv.Itoa(numeroInodo) + ":" + strconv.Itoa(Contador+1) + " -> Bloque" + strconv.Itoa(int(i)) + ":0;\n"
			Dot += "Bloque" + strconv.Itoa(int(i)) + "[label = <\n"
			Dot += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
			DesplazamientoBloque := int(sb.S_block_start) + (int(i) * binary.Size(FolderBlock{}))
			carpeta := FolderBlock{}
			archivo.Seek(int64(DesplazamientoBloque), 0)
			binary.Read(archivo, binary.LittleEndian, &carpeta)
			if inodo.I_type == [1]byte{'0'} {
				Dot += "<tr><td colspan=\"2\" port='0'>Bloque" + strconv.Itoa(int(i)) + "</td></tr>\n"
				Contador2 := 0
				for _, j := range carpeta.B_content {
					fmt.Println("Nombre: ", string(j.B_name[:]))
					nam := strings.TrimRight(string(j.B_name[:]), string(rune(0)))

					if Contador2 == 0 {
						nam = "."
					}
					if Contador2 == 1 {
						nam = ".."
					}
					if j.B_inodo == -1 {
						nam = ""
					}
					fmt.Println("Nombre: ", nam)
					Dot += "<tr><td>" + nam + "</td><td port='" + strconv.Itoa(Contador2+1) + "'>" + strconv.Itoa(int(j.B_inodo)) + "</td></tr>\n"
					Contador2++
				}
				Dot += "</table>>];\n"
				Contador2 = 0
				for _, j := range carpeta.B_content {
					if j.B_inodo != -1 {
						if j.B_name[0] != '.' {
							//Leer el inodo
							Dot += "Bloque" + strconv.Itoa(int(i)) + ":" + strconv.Itoa(Contador2+1) + " -> Inodo" + strconv.Itoa(int(j.B_inodo)) + ":0;\n"
							//Buscar el inodo siguiente
							DesplazamientoInodo := int(sb.S_inode_start) + (int(j.B_inodo) * binary.Size(Inodes{}))
							inodoSiguiente := NewInodes()
							archivo.Seek(int64(DesplazamientoInodo), 0)
							binary.Read(archivo, binary.LittleEndian, &inodoSiguiente)
							Dot += RecursivoTree(inodoSiguiente, sb, archivo, int(j.B_inodo))
						}
					}
					Contador2++
				}
			} else {
				file := Fileblock{}
				archivo.Seek(int64(DesplazamientoBloque), 0)
				binary.Read(archivo, binary.LittleEndian, &file)
				Dot += "<tr><td colspan=\"1\" port='0'>Bloque" + strconv.Itoa(int(i)) + "</td></tr>\n"
				Dot += "<tr><td port='1'>" + strings.TrimRight(string(file.B_content[:]), string(rune(0))) + "</td></tr>\n"
				Dot += "</table>>];\n"
			}
		}
		Contador++
	}

	return Dot
}

func ReporteSB(idValor string, pathValor string) {
	//Abrir el disco A
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	//Buscar la particion montada con el ID
	Particion := VerificarParticionMontada(idValor)
	if Particion == -1 {
		fmt.Println("No se encontro la particion montada con el ID: ", idValor)
		return
	}
	MountActual := particionesMontadas[Particion]
	//Abrir el disco
	archivo, err := os.OpenFile(""+MountActual.Path+"", os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	//Leer el superbloque
	sb := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sb)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}
	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte SuperBlock \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	//se le agrega el nombre del disco segun la ruta a la tabla
	Dot += "<TR><TD bgcolor=\"lightgrey\">Disco</TD><TD>" + MountActual.Path + "</TD></TR>\n"

	Dot += "<TR><TD bgcolor=\"lightgrey\">s_filesystem_type</TD><TD>" + strconv.Itoa(int(sb.S_filesystem_type)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inodes_count</TD><TD>" + strconv.Itoa(int(sb.S_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_blocks_count</TD><TD>" + strconv.Itoa(int(sb.S_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_blocks_count</TD><TD>" + strconv.Itoa(int(sb.S_free_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_inodes_count</TD><TD>" + strconv.Itoa(int(sb.S_free_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mtime</TD><TD>" + string(sb.S_mtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_umtime</TD><TD>" + string(sb.S_umtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mnt_count</TD><TD>" + strconv.Itoa(int(sb.S_mnt_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_magic</TD><TD>" + strconv.Itoa(int(sb.S_magic)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_size</TD><TD>" + strconv.Itoa(int(sb.S_inode_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_size</TD><TD>" + strconv.Itoa(int(sb.S_block_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_ino</TD><TD>" + strconv.Itoa(int(sb.S_first_ino)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_blo</TD><TD>" + strconv.Itoa(int(sb.S_first_blo)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_inode_start</TD><TD>" + strconv.Itoa(int(sb.S_bm_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_block_start</TD><TD>" + strconv.Itoa(int(sb.S_bm_block_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_start</TD><TD>" + strconv.Itoa(int(sb.S_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_start</TD><TD>" + strconv.Itoa(int(sb.S_block_start)) + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"

	//Quitar la extension al archivo (pdf, etc, )
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
	archivoDot, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		return
	}
	defer archivoDot.Close()
	_, err = archivoDot.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		return
	}

	//Generar la imagen
	//Quitar el punto
	extensionSinPunto := strings.TrimPrefix(extension, ".")
	//Correr con todos los permisos
	cmd := exec.Command("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	fmt.Println("dot", "-T", extensionSinPunto, DotName, "-o", dirPath+fileName+extension)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)

		return
	}

	fmt.Println("Reporte sb generado con exito")
}

func BM_inode(idValor string, pathValor string) {

}

func ReporteMBR(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	// Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	// Buscar la partición montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La partición no está montada\n"
		return respuesta
	}

	MountActual := particionesMontadas[indice]

	// Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()

	// Leer el MBR
	disk := MBR{}
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		respuesta += "Error al leer el MBR\n"
		fmt.Println("Error al leer el MBR")
		return respuesta
	}

	// Crear el contenido del archivo .dot
	Dot := "digraph G {bgcolor=\"slategrey\" label=\"Reporte MBR\" layout=dot "
	Dot += "labelloc = \"t\" edge [weight=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"] "
	Dot += "node[shape=record, color=lightgrey] nodoMBR[label=\"MBR | { mbr_tamanio | " + fmt.Sprintf("%d", disk.Mbr_tamano) + " } |"
	Dot += "{ mbr_fecha_creacion | " + string(disk.Mbr_fecha_creacion[:]) + " } |"
	Dot += "{ mbr_disk_signature | " + fmt.Sprintf("%d", disk.Mbr_disk_signature) + " } |"

	// Agregar información de las particiones
	// Agregar información de las particiones
	for i := 1; i <= 4; i++ {
		particion := obtenerParticion(disk, i)
		if particion.Part_size != 0 {
			Dot += "| { Particion " + fmt.Sprintf("%d", i) + " | { part_status | " + string(particion.Part_status[:]) + " } |"
			Dot += "{ part_type | " + string(particion.Part_type[:]) + " } |"
			Dot += "{ part_fit | " + string(particion.Part_fit[:]) + " } |"
			Dot += "{ part_start | " + fmt.Sprintf("%d", particion.Part_start) + " } |"
			Dot += "{ part_size | " + fmt.Sprintf("%d", particion.Part_size) + " } |"
			Dot += "{ part_name | " + strings.Trim(string(particion.Part_name[:]), "\x00") + " } }" // Asegurarse de limpiar la cadena
		}
	}

	Dot += "\"];\n}"

	// Crear el archivo .dot
	extension := path.Ext(pathValor)
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"

	file, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot")
		respuesta += "Error al crear el archivo .dot\n"
		return respuesta
	}
	defer file.Close()

	// Escribir el contenido en el archivo .dot
	fmt.Println("Contenido de Dot generado:\n" + Dot)

	_, err = file.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot")
		respuesta += "Error al escribir el archivo .dot\n"
		return respuesta
	}

	fmt.Println("Archivo .dot creado")

	// Quitar el punto a la extensión
	extension = extension[1:]

	// Crear el reporte utilizando Graphviz
	cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
	fmt.Println("Ejecutando comando: dot -T" + extension + " " + DotName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al crear el reporte:", err)
		respuesta += "Error al crear el reporte\n"
		return respuesta
	}

	return "Reporte MBR creado con éxito\n"
}

// Función para obtener una partición del MBR según el índice
func obtenerParticion(mbr MBR, index int) Partition {
	switch index {
	case 1:
		return mbr.Mbr_partition_1
	case 2:
		return mbr.Mbr_partition_2
	case 3:
		return mbr.Mbr_partition_3
	case 4:
		return mbr.Mbr_partition_4
	default:
		return Partition{}
	}
}

func ReporteBMInode(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	//Leer el bitmap de inodos, teniendo 20 registros por fila
	Desplazamiento := int(superBloque.S_bm_inode_start)
	BmString := ""

	for i := 0; i < int(superBloque.S_inodes_count); i++ {
		var bit byte
		archivo.Seek(int64(Desplazamiento+i), 0)
		err = binary.Read(archivo, binary.LittleEndian, &bit)
		if err != nil {
			respuesta += "Error al leer el bitmap de inodos\n"
			fmt.Println("Error al leer el bitmap de inodos")
			return respuesta
		}
		if bit == 0 {
			BmString += "0"
		} else {
			BmString += "1"
		}
		if (i+1)%20 == 0 {
			BmString += "\n"
		}
	}
	Dot := "digraph G{\n"
	Dot += "a[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Bitmap de Inodos</TD></TR>\n"
	Dot += "<TR><TD>" + BmString + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"
	//Crear el archivo dot
	extension := path.Ext(pathValor)
	if extension == ".txt" {
		//Crear el archivo .txt
		file, err := os.Create(pathValor)
		if err != nil {
			fmt.Println("Error al crear el archivo .txt")
			respuesta += "Error al crear el archivo .txt\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .txt
		_, err = file.WriteString(BmString)
		if err != nil {
			fmt.Println("Error al escribir el archivo .txt")
			respuesta += "Error al escribir el archivo .txt\n"
			return respuesta
		}
		fmt.Println("Archivo .txt creado")
		return "Reporte Bitmap de Inodos creado con exito\n"

	} else {
		//Archivo sin extension
		fileName = strings.TrimSuffix(fileName, extension)
		DotName := dirPath + fileName + ".dot"
		//Crear el archivo .dot
		file, err := os.Create(DotName)
		if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .dot
		_, err = file.WriteString(Dot)
		if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
		}
		fmt.Println("Archivo .dot creado")

		//Quitar el punto a la extension
		extension = extension[1:]

		//Crear el reporte
		cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
		fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error al crear el reporte")
			respuesta += "Error al crear el reporte\n"
			return respuesta
		}

		return "Reporte Bitmap de Inodos creado con exito\n"
	}
}

func ReporteBMBlock(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	//Leer el bitmap de bloques, teniendo 20 registros por fila
	Desplazamiento := int(superBloque.S_bm_block_start)
	BmString := ""

	for i := 0; i < int(superBloque.S_blocks_count); i++ {
		var bit byte
		archivo.Seek(int64(Desplazamiento+i), 0)
		err = binary.Read(archivo, binary.LittleEndian, &bit)
		if err != nil {
			respuesta += "Error al leer el bitmap de bloques\n"
			fmt.Println("Error al leer el bitmap de bloques")
			return respuesta
		}
		if bit == 0 {
			BmString += "0"
		} else {
			BmString += "1"
		}
		if (i+1)%20 == 0 {
			BmString += "\n"
		}
	}
	Dot := "digraph G{\n"
	Dot += "a[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Bitmap de Bloques</TD></TR>\n"
	Dot += "<TR><TD>" + BmString + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"
	//Crear el archivo dot
	extension := path.Ext(pathValor)
	if extension == ".txt" {
		//Crear el archivo .txt
		file, err := os.Create(pathValor)
		if err != nil {
			fmt.Println("Error al crear el archivo .txt")
			respuesta += "Error al crear el archivo .txt\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .txt
		_, err = file.WriteString(BmString)
		if err != nil {
			fmt.Println("Error al escribir el archivo .txt")
			respuesta += "Error al escribir el archivo .txt\n"
			return respuesta
		}
		fmt.Println("Archivo .txt creado")
		return "Reporte Bitmap de Bloques creado con exito\n"

	} else {
		//Archivo sin extension
		fileName = strings.TrimSuffix(fileName, extension)
		DotName := dirPath + fileName + ".dot"
		//Crear el archivo .dot
		file, err := os.Create(DotName)
		if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .dot
		_, err = file.WriteString(Dot)
		if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
		}
		fmt.Println("Archivo .dot creado")

		//Quitar el punto a la extension
		extension = extension[1:]

		//Crear el reporte
		cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
		fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error al crear el reporte")
			respuesta += "Error al crear el reporte\n"
			return respuesta
		}

		return "Reporte Bitmap de Bloques creado con exito\n"
	}
}
