insert into Usuarios values ('Prueba', 'prueba');
insert into Usuarios values ('Admin', 'admin');
insert into Usuarios values ('Socia', 'Socia');

INSERT INTO `inmosoftDB`.`Tipo_Propiedad` values (1, 'casa');
INSERT INTO `inmosoftDB`.`Tipo_Propiedad` values (2, 'bodega');
INSERT INTO `inmosoftDB`.`Tipo_Propiedad` values (3, 'local');
INSERT INTO `inmosoftDB`.`Tipo_Propiedad` values (4, 'apartamento');
INSERT INTO `inmosoftDB`.`Tipo_Propiedad` values (5, 'terreno');

select * from Tipo_Propiedad;

INSERT INTO `inmosoftDB`.`Propietario` 
(`id_propietario`, `nombre_propietario`, `apellido_paterno_propietario`, `apellido_materno_propietario`, `telefono_propietario`, `correo_propietario`) 
VALUES 
(1, 'Carlos', 'Garcia', 'Lopez', '555-1234', 'carlos.garcia@example.com');

INSERT INTO `inmosoftDB`.`Propietario` 
(`id_propietario`, `nombre_propietario`, `apellido_paterno_propietario`, `apellido_materno_propietario`, `telefono_propietario`, `correo_propietario`) 
VALUES 
(2, 'Maria', 'Fernandez', 'Rodriguez', '555-5678', 'maria.fernandez@example.com');

INSERT INTO `inmosoftDB`.`Propietario` 
(`id_propietario`, `nombre_propietario`, `apellido_paterno_propietario`, `apellido_materno_propietario`, `telefono_propietario`, `correo_propietario`) 
VALUES 
(3, 'Luis', 'Perez', 'Martinez', '555-9876', 'luis.perez@example.com');

select * from Propietario;

-- Insert 1
INSERT INTO `inmosoftDB`.`Propiedades` 
(`id_propiedad`, `titulo`, `fecha_alta`, `direccion`, `colonia`, `ciudad`, `referencia`, `precio`, 
 `mts_construccion`, `mts_terreno`, `habitada`, `amueblada`, `num_plantas`, 
 `num_recamaras`, `num_banos`, `size_cochera`, `mts_jardin`, `gas`, `comodidades`, 
 `extras`, `utilidades`, `observaciones`, `id_tipo_propiedad`, `id_propietario`, `usuario`)
VALUES 
(1, 'Casa Sorrento Residencial', '2024-10-11', 'Los Pastores S/N, Los Valdez, Fraccionamiento Sorrento, Saltillo', 'Fraccionamiento Sorrento', 'Saltillo', 'Frente al parque', 3463000, 
 144, 173, 1, 1, 2, 4, 3, 2, 60, 'natural', 'clima,aljibe', 
 'alberca,jardin,techada', 'agua,luz,internet', 
 'El Modelo Bari II es la eleccion perfecta para quienes buscan una combinacion de estilo contemporaneo y funcionalidad en su nuevo hogar. Disenado para maximizar el confort y la comodidad, este modelo se adapta a las necesidades de familias modernas.', 1, 1, 'Admin');

-- Insert 2
INSERT INTO `inmosoftDB`.`Propiedades` 
(`id_propiedad`, `titulo`, `fecha_alta`, `direccion`, `colonia`, `ciudad`, `referencia`, `precio`, 
 `mts_construccion`, `mts_terreno`, `habitada`, `amueblada`, `num_plantas`, 
 `num_recamaras`, `num_banos`, `size_cochera`, `mts_jardin`, `gas`, `comodidades`, 
 `extras`, `utilidades`, `observaciones`, `id_tipo_propiedad`, `id_propietario`, `usuario`)
VALUES 
(2, 'Casa Nueva en Privadas Mirasierra Blvd Mirasierra', '2024-10-12', 'BLVD CEDROS 3ERA ETAPA NUEVA MIRASIERRA, Nazario S Ortiz Garza, Saltillo', 'Mirasierra', 'Saltillo', 'Cerca de la escuela', 2383000, 
 150, 300, 0, 0, 1, 2, 1, 1, 30, 'estacionario', 'clima,calefaccion', 
 'techada,cocineta', 'agua,luz', 
 'Casas totalmente nuevas en Privanzas MiraSierra Cerca de todo lo que necesitas EN TOTAL tenemos 4 modelos Desde $2,384,000 hasta $2,550,000 Se acepta COFINAVIT E INFONAVIT y credito bancario. Muy exclusivo y privado con caseta y acceso controlado.', 1, 2, 'Admin');

-- Insert 3
INSERT INTO `inmosoftDB`.`Propiedades` 
(`id_propiedad`, `titulo`, `fecha_alta`, `direccion`, `colonia`, `ciudad`, `referencia`, `precio`, 
 `mts_construccion`, `mts_terreno`, `habitada`, `amueblada`, `num_plantas`, 
 `num_recamaras`, `num_banos`, `size_cochera`, `mts_jardin`, `gas`, `comodidades`, 
 `extras`, `utilidades`, `observaciones`, `id_tipo_propiedad`, `id_propietario`, `usuario`)
VALUES 
(3, 'Casa de una planta', '2024-10-13', 'Residencial San Alberto, Residencial San Alberto, Saltillo', 'San Alberto', 'Saltillo', 'Cerca del H-E-B', 10590000, 
 180, 350, 1, 0, 2, 3, 2, 2, 50, 'natural', 'calefaccion,hidroneumatico', 
 'jardin,cuarto_servicio', 'agua,luz,internet', 
 'Casa de una planta en venta en Residencial San Alberto Al norte de la Ciudad de Saltillo, Coahuila. Fraccionamiento muy seguro y privado. Recibidor, 1/2 bano de visitas, Sala, comedor, cocina integral con barra desayunador con cubierta de granito, amplia terraza exterior con asador y otro medio bano, bodega exterior, 2 recamaras muy amplias con bano - vestidor compartido (suficientemente amplio para adaptarse otro bano) una de las recamaras con closet, lavanderia, bodega interior y otra exterior, jardin en la cochera que es techada para 2 autos y pasillo de servicio.', 1, 3, 'Admin');

-- Insert 4
INSERT INTO `inmosoftDB`.`Propiedades` 
(`id_propiedad`, `titulo`, `fecha_alta`, `direccion`, `colonia`, `ciudad`, `referencia`, `precio`, 
 `mts_construccion`, `mts_terreno`, `habitada`, `amueblada`, `num_plantas`, 
 `num_recamaras`, `num_banos`, `size_cochera`, `mts_jardin`, `gas`, `comodidades`, 
 `extras`, `utilidades`, `observaciones`, `id_tipo_propiedad`, `id_propietario`, `usuario`)
VALUES 
(4, 'Casa en Fraccionamiento Nuestra Sra de Fatima', '2024-10-14', 'Fracc. Nuestra Sra de Fatima, Nuestra Senora de Fatima, Saltillo', 'Nuestra Senora de Fatima', 'Saltillo', 'Junto al mercado', 4520000, 
 170, 320, 0, 1, 1, 2, 1, 1, 40, 'estacionario', 'clima', 
 'techada', 'agua', 
 'Casa en venta en Fracc Nuestra Senora de Fatima Fraccionamiento cerrado con ubicacion privilegiada, cercano a centros comerciales, escuelas, universidades y hospitales, facil acceso a Carretera Saltillo-Monterrey', 2, 1, 'Admin');

-- Insert 5
INSERT INTO `inmosoftDB`.`Propiedades` 
(`id_propiedad`, `titulo`, `fecha_alta`, `direccion`, `colonia`, `ciudad`, `referencia`, `precio`, 
 `mts_construccion`, `mts_terreno`, `habitada`, `amueblada`, `num_plantas`, 
 `num_recamaras`, `num_banos`, `size_cochera`, `mts_jardin`, `gas`, `comodidades`, 
 `extras`, `utilidades`, `observaciones`, `id_tipo_propiedad`, `id_propietario`, `usuario`)
VALUES 
(5, 'Casa en Venta en Col. Jardines del Lago Saltillo', '2024-10-15', 'Col. Jardines del LAGO, Jardines del Lago, Saltillo', 'Jardines de Lago', 'Saltillo', 'Junto a la iglesia', 4700000, 
 200, 400, 1, 0, 2, 4, 3, 3, 80, 'natural', 'clima,calefaccion', 
 'alberca,jardin', 'agua,luz,internet', 
 'Disfruta de la comodidad y el lujo en esta casa de un solo piso en Col. Jardines de Lago. Con 328m2 de terreno y 300m2 de construccion, cuenta con 3 recamaras, 2.5 banos, area de juegos con chimenea, cocina equipada, cuarto de servicio, cochera techada y fraccionamiento privado. ¡No pierdas esta oportunidad! Precio: $4,700,000.00. #CasaEnVenta #JardinesDelLago #Saltillo #CasaDeUnPiso #FraccionamientoPrivado', 1, 2, 'Admin');

select * from Propiedades;

INSERT INTO `inmosoftDB`.`Estado_Propiedades` values (1, 'venta', 'disponible', null, 1);
INSERT INTO `inmosoftDB`.`Estado_Propiedades` values (2, 'renta', 'rentada', null, 2);
INSERT INTO `inmosoftDB`.`Estado_Propiedades` values (3, 'venta', 'vendida', null, 3);
INSERT INTO `inmosoftDB`.`Estado_Propiedades` values (4, 'renta', 'disponible', null, 4);
INSERT INTO `inmosoftDB`.`Estado_Propiedades` values (5, 'venta', 'disponible', null, 5);

select * from Estado_Propiedades;

INSERT INTO `inmosoftDB`.`Prospecto` 
(`id_cliente`, `nombre_prospecto`, `apellido_paterno_prospecto`, `apellido_materno_prospecto`, `telefono_prospecto`, `correo_prospecto`) 
VALUES 
(1, 'Alejandra', 'Martinez', 'Perez', '555-1234', 'ale.mtz@example.com'),
(2, 'Mario', 'Lopez', 'Martinez', '555-5678', 'mario.lopez@example.com'),
(3, 'Carlos', 'Hernandez', 'Sanchez', '555-8765', 'carlos.hernandez@example.com'),
(4, 'Ana', 'Ramirez', 'Torres', '555-4321', 'ana.ramirez@example.com'),
(5, 'Deyanira', 'Hernandez', 'Hernandez', '555-2468', 'deyanirahdzhdz@example.com'),
(6, 'Elena', 'Gutierrez', 'Diaz', '555-1357', 'elena.gutierrez@example.com'),
(7, 'Josefa', 'Mendoza', 'Fernandez', '555-9876', 'josefa.mendoza@example.com');

INSERT INTO `inmosoftDB`.`Citas` 
(`id_citas`, `titulo_cita`, `fecha_cita`, `hora_cita`, `descripcion_cita`, `usuario`, `id_cliente`) 
VALUES 
(1, 'Primera visita de inspeccion', '2024-12-12', 1700, 'Primera inspeccion para revisar el estado de la propiedad.', 'Admin', 1),
(2, 'Negociacion', '2024-12-12', 1000, 'Seguimiento a observaciones de la primera inspeccion.', 'Admin', 2),
(3, 'Revision de documentos', '2024-12-15', 1500, 'Revision de papelería necesaria para el proceso.', 'Admin', 3),
(4, 'Confirmacion de oferta', '2024-12-15', 1100, 'Confirmar la oferta presentada por el cliente.', 'Admin', 4),
(5, 'Visita para negociacion', '2024-12-15', 1400, 'Negociacion de condiciones finales.', 'Admin', 5),
(6, 'Revision final de contrato', '2024-12-15', 1600, 'Revisar el contrato antes de la firma final.', 'Admin', 6),
(7, 'Firma de contrato', '2024-12-18', 1200, 'Firma oficial del contrato con el cliente.', 'Admin', 7);

INSERT INTO `inmosoftDB`.`ImagenesProspecto`
(`id_imagen`, `ruta_imagen`, `descripcion_imagen`, `principal`, `id_prospecto`)
VALUES
(1, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto1.jpg', "null", 1, 1),
(2, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto2.jpg', "null", 1, 2),
(3, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto3.jpg', "null", 1, 3),
(4, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto4.jpg', "null", 1, 4),
(5, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto5.jpg', "null", 1, 5),
(6, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto6.jpg', "null", 1, 6),
(7, 'C:\\Users\\josel\\Documents\\Imagenes\\Prospectos\\prospecto7.jpg', "null", 1, 7);


INSERT INTO `inmosoftDB`.`Imagenes` (`id_imagen`, `ruta_imagen`, `descripcion_imagen`, `principal`, `id_propiedad`) 
VALUES 

(1, 'C:\\Users\\josel\\Documents\\Imagenes\\casa1-1.jpg', "null", 1, 1),
(2, 'C:\\Users\\josel\\Documents\\Imagenes\\casa1-2.jpg', "null", 0, 1),
(3, 'C:\\Users\\josel\\Documents\\Imagenes\\casa1-3.jpg', "null", 0, 1),
(4, 'C:\\Users\\josel\\Documents\\Imagenes\\casa1-4.jpg', "null", 0, 1),

(5, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-1.jpg', "null", 1, 2),
(6, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-2.jpg', "null", 0, 2),
(7, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-3.jpg', "null", 0, 2),
(8, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-4.jpg', "null", 0, 2),
(9, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-5.jpg', "null", 0, 2),
(10, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-6.jpg', "null", 0, 2),
(11, 'C:\\Users\\josel\\Documents\\Imagenes\\casa2-7.jpg', "null", 0, 2),

(12, 'C:\\Users\\josel\\Documents\\Imagenes\\casa3-1.jpg', "null", 1, 3),
(13, 'C:\\Users\\josel\\Documents\\Imagenes\\casa3-2.jpg', "null", 0, 3),
(14, 'C:\\Users\\josel\\Documents\\Imagenes\\casa3-3.jpg', "null", 0, 3),
(15, 'C:\\Users\\josel\\Documents\\Imagenes\\casa3-4.jpg', "null", 0, 3),
(16, 'C:\\Users\\josel\\Documents\\Imagenes\\casa3-5.jpg', "null", 0, 3),
(17, 'C:\\Users\\josel\\Documents\\Imagenes\\casa3-6.jpg', "null", 0, 3),

(18, 'C:\\Users\\josel\\Documents\\Imagenes\\casa4-1.jpg', "null", 1, 4),
(19, 'C:\\Users\\josel\\Documents\\Imagenes\\casa4-2.jpg', "null", 0, 4),
(20, 'C:\\Users\\josel\\Documents\\Imagenes\\casa4-3.jpg', "null", 0, 4),
(21, 'C:\\Users\\josel\\Documents\\Imagenes\\casa4-4.jpg', "null", 0, 4),
(22, 'C:\\Users\\josel\\Documents\\Imagenes\\casa4-5.jpg', "null", 0, 4),
(23, 'C:\\Users\\josel\\Documents\\Imagenes\\casa4-6.jpg', "null", 0, 4),

(24, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-1.jpg', "null", 1, 5),
(25, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-2.jpg', "null", 0, 5),
(26, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-3.jpg', "null", 0, 5),
(27, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-4.jpg', "null", 0, 5),
(28, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-5.jpg', "null", 0, 5),
(29, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-6.jpg', "null", 0, 5),
(30, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-7.jpg', "null", 0, 5),
(31, 'C:\\Users\\josel\\Documents\\Imagenes\\casa5-8.jpg', "null", 0, 5);

INSERT INTO `inmosoftDB`.`Contratos` 
(`id_contrato`, `titulo_contrato`, `descripcion_contrato`, `tipo`, `ruta_pdf`, `id_propiedad`) 
VALUES
(1, 'Contrato de Arrendamiento', 'Documento legal que detalla los terminos y condiciones del arrendamiento de la propiedad, incluyendo plazos, costos y responsabilidades de ambas partes.', 'Arrendamiento', 'C:\\Users\\josel\\Documents\\Imagenes\\Contratos\\contrato1.pdf', 1),
(2, 'Contrato de Compra-Venta', 'Acuerdo formal que establece las condiciones para la transferencia de propiedad entre el comprador y el vendedor, especificando precio, plazos y obligaciones legales.', 'Compraventa', 'C:\\Users\\josel\\Documents\\Imagenes\\Contratos\\contrato2.pdf', 2),
(3, 'Contrato de Maria Casa Norte', 'Contrato que regula el arrendamiento de una propiedad ubicada en el area norte, incluyendo las clausulas de pago y uso adecuado del inmueble.', 'Arrendamiento', 'C:\\Users\\josel\\Documents\\Imagenes\\Contratos\\contrato3.pdf', 3),
(4, 'Contrato de Juan de Enero', 'Documento legal para la renta de una propiedad, estableciendo los derechos y responsabilidades del arrendador y arrendatario durante el periodo acordado.', 'Arrendamiento', 'C:\\Users\\josel\\Documents\\Imagenes\\Contratos\\contrato4.pdf', 4),
(5, 'Contrato de Renta Casa Fundadores', 'Acuerdo que detalla los terminos para el arrendamiento de una propiedad en la zona de Fundadores, con estipulaciones sobre mantenimiento, duracion y pagos.', 'Arrendamiento', 'C:\\Users\\josel\\Documents\\Imagenes\\Contratos\\contrato5.pdf', 5);
