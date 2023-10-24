# Grupo13-Laboratorio-1

| Nombre Integrante | Rol |
| ------------ | ------------ | 
| Diego Acevedo | 202073532-8 | 
| Vicente Ruiz | 202073585-9  | 
| Gabriel Vergara | 202073616-2 | 


# Instrucciones
Orden sugerido de comandos:

Primero se debe correr el siguiente comando en cualquiera de las máquina (ojalá el número de máquina que sale en el comando):

```
sudo make docker-vm049
```
```
sudo make docker-vm050
```

```
sudo make docker-vm051
```
```
sudo make docker-vm052
```

Donde cada VM da inicio a los procesos de:

- vm049: OMS y Australia
- vm050: Onu y Europa
- vm051: DataNode1 y Asia
- vm052: DataNode2 y Latinoamerica

 
# Consideraciones
El input de la terminal que usa la ONU no está funcional, está implementado que se _hardcodea_ la condición a consultar y luego se comunica de correcta manera a OMS para luego obtener los nombres de los DataNode correspondientes

