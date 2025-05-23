# Prueba Técnica - Tienda de Deporte

API REST en Golang para gestionar productos deportivos.  
Incluye CRUD, métricas.

## Variables de entorno

Se debe renombrar env.txt a .env para que tome las variables de entorno

## Base de datos Seleccionada MongoDB

# Justificación

se selecciona MongoDB teniendo en cuenta rendimiento escalabilidad tem de concurrencia mejor integracion si cmabian los esquemas o nuevas estructuras y mayor control de los datos

```bash
  MONGO_URI=
  MONGO_DB_NAME=
```

## Para ejecutar el proyecto

```bash
  make  run
  - Este comando instala lo necesario ejecuta las pruebas unitarias y inicia el servicio el cual esta en :8080
```

# deploy en railway

![Logo de MLSport](deploy.png)
https://mlsport-production.up.railway.app/

## Documentación de la API

![Logo de MLSport](api.png)
http://localhost:8080/api

## CRUD de Productos Deportivos

Estos se encuentran en la ruta principal /

## Métricas de Productos

se encuentran en el endpoint /api/products/dashboard
