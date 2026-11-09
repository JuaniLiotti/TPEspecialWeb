// =====================================================================
// Italpiel - Frontend (TP4)
// Consume la API REST del servidor Go con fetch y arma la vista en el DOM.
// =====================================================================

const API_CUEROS = '/api/cueros/';
const API_PRODUCTOS = '/api/products/';

const formatoPrecio = new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
});


// ---------------------------------------------------------------------
// Utilidades de comunicación con la API
// ---------------------------------------------------------------------

// Si la respuesta no es 2xx, lanza un error con el mensaje que mandó el servidor
// (los handlers de Go responden los errores como texto con http.Error)
async function verificarRespuesta(respuesta) {
    if (!respuesta.ok) {
        const mensaje = (await respuesta.text()).trim();
        throw new Error(mensaje || `Error ${respuesta.status}`);
    }
    return respuesta;
}

async function obtenerLista(url) {
    const respuesta = await verificarRespuesta(await fetch(url));
    return respuesta.json();
}

// Crea la entidad con un POST JSON y, si se eligieron archivos, los sube
// en un segundo POST multipart/form-data a /api/.../{id}/imagenes
async function crearConImagenes(url, datos, archivos) {
    const respuesta = await verificarRespuesta(await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(datos),
    }));
    const creado = await respuesta.json();

    if (archivos.length > 0) {
        const formData = new FormData();
        for (const archivo of archivos) {
            formData.append('imagenes', archivo); // mismo nombre de campo que espera el servidor
        }
        // No se setea Content-Type: el navegador lo arma solo, con el "boundary" del multipart
        const respuestaImagenes = await fetch(`${url}${creado.id}/imagenes`, {
            method: 'POST',
            body: formData,
        });
        if (!respuestaImagenes.ok) {
            const mensaje = (await respuestaImagenes.text()).trim();
            throw new Error(`Se guardó "${creado.nombre}", pero no se pudieron subir las imágenes: ${mensaje}`);
        }
    }
    return creado;
}

async function eliminar(url, id) {
    await verificarRespuesta(await fetch(`${url}${id}`, { method: 'DELETE' }));
}


// ---------------------------------------------------------------------
// Construcción del DOM
// ---------------------------------------------------------------------

// Atajo para crear un elemento con clase y texto.
// Usa textContent (no innerHTML) para que los datos se muestren como texto
// y nunca se interpreten como HTML.
function crearElemento(etiqueta, clase, texto) {
    const elemento = document.createElement(etiqueta);
    if (clase) elemento.className = clase;
    if (texto !== undefined) elemento.textContent = texto;
    return elemento;
}

// Arma una tarjeta genérica. "detalles" es una lista de textos que cambia
// según sea un cuero o un producto.
function crearTarjeta(item, detalles) {
    const tarjeta = crearElemento('article', 'tarjeta');
    tarjeta.dataset.id = item.id;

    // Portada: la primera imagen, o un recuadro si no tiene ninguna
    if (item.imagenes.length > 0) {
        const img = crearElemento('img', 'tarjeta-imagen');
        img.src = item.imagenes[0].url;
        img.alt = item.nombre;
        tarjeta.append(img);
    } else {
        tarjeta.append(crearElemento('div', 'tarjeta-imagen sin-imagen', 'Sin imagen'));
    }

    const cuerpo = crearElemento('div', 'tarjeta-cuerpo');
    cuerpo.append(crearElemento('h4', 'tarjeta-titulo', item.nombre));

    for (const detalle of detalles) {
        if (detalle) cuerpo.append(crearElemento('p', 'tarjeta-detalle', detalle));
    }
    if (item.imagenes.length > 1) {
        cuerpo.append(crearElemento('p', 'tarjeta-detalle', `${item.imagenes.length} fotos`));
    }

    const pie = crearElemento('div', 'tarjeta-pie');
    pie.append(crearElemento('span', 'tarjeta-precio', formatoPrecio.format(item.precio)));
    const botonEliminar = crearElemento('button', 'btn-eliminar', 'Eliminar');
    botonEliminar.type = 'button';
    botonEliminar.dataset.id = item.id; // el listener delegado lee este ID
    pie.append(botonEliminar);

    cuerpo.append(pie);
    tarjeta.append(cuerpo);
    return tarjeta;
}

// Reemplaza el contenido de un contenedor por las tarjetas (o un mensaje si no hay)
function mostrarLista(contenedor, items, armarDetalles, mensajeVacio) {
    contenedor.replaceChildren();

    if (items.length === 0) {
        contenedor.append(crearElemento('p', 'lista-vacia', mensajeVacio));
        return;
    }
    for (const item of items) {
        contenedor.append(crearTarjeta(item, armarDetalles(item)));
    }
}

function mostrarErrorEnLista(contenedor, error) {
    contenedor.replaceChildren(crearElemento('p', 'mensaje-error', `No se pudo cargar la lista: ${error.message}`));
}


// ---------------------------------------------------------------------
// Cueros
// ---------------------------------------------------------------------

const listaCueros = document.getElementById('lista-cueros');
const formCuero = document.getElementById('form-cuero');

async function cargarCueros() {
    try {
        const cueros = await obtenerLista(API_CUEROS);
        mostrarLista(listaCueros, cueros, (c) => [
            `${c.tipo}${c.acabado ? ' · ' + c.acabado : ''}`,
            `${c.color ? c.color + ' · ' : ''}${c.espesor_mm} mm`,
            c.descripcion,
        ], 'Todavía no hay cueros cargados.');
    } catch (error) {
        mostrarErrorEnLista(listaCueros, error);
    }
}

formCuero.addEventListener('submit', async (evento) => {
    evento.preventDefault(); // evita que el navegador envíe el form y recargue la página

    const f = formCuero.elements;
    const datos = {
        nombre: f.nombre.value.trim(),
        tipo: f.tipo.value.trim(),
        acabado: f.acabado.value.trim(),
        color: f.color.value.trim(),
        espesor_mm: parseFloat(f.espesor_mm.value),
        descripcion: f.descripcion.value.trim(),
        precio: parseInt(f.precio.value, 10),
    };

    await enviarFormulario(formCuero, () => crearConImagenes(API_CUEROS, datos, f.imagenes.files), cargarCueros);
});


// ---------------------------------------------------------------------
// Productos terminados
// ---------------------------------------------------------------------

const listaProductos = document.getElementById('lista-productos');
const formProducto = document.getElementById('form-producto');

async function cargarProductos() {
    try {
        const productos = await obtenerLista(API_PRODUCTOS);
        mostrarLista(listaProductos, productos, (p) => [
            [p.categoria, p.color].filter(Boolean).join(' · '),
            p.descripcion,
        ], 'Todavía no hay productos cargados.');
    } catch (error) {
        mostrarErrorEnLista(listaProductos, error);
    }
}

formProducto.addEventListener('submit', async (evento) => {
    evento.preventDefault();

    const f = formProducto.elements;
    const datos = {
        nombre: f.nombre.value.trim(),
        categoria: f.categoria.value.trim(),
        color: f.color.value.trim(),
        descripcion: f.descripcion.value.trim(),
        precio: parseInt(f.precio.value, 10),
    };

    await enviarFormulario(formProducto, () => crearConImagenes(API_PRODUCTOS, datos, f.imagenes.files), cargarProductos);
});


// ---------------------------------------------------------------------
// Lógica compartida de formularios y borrado
// ---------------------------------------------------------------------

// Envía un formulario: deshabilita el botón mientras espera, muestra el error
// del servidor si lo hay y, al terminar, recarga la lista.
async function enviarFormulario(form, accion, recargarLista) {
    const boton = form.querySelector('button[type="submit"]');
    const mensajeError = form.querySelector('.mensaje-error');

    boton.disabled = true;
    mensajeError.hidden = true;
    try {
        await accion();
        form.reset();
    } catch (error) {
        mensajeError.textContent = error.message;
        mensajeError.hidden = false;
    } finally {
        boton.disabled = false;
        // Aunque falle la subida de imágenes, la entidad pudo haberse creado
        await recargarLista();
    }
}

// Delegación de eventos: un solo listener en el contenedor atiende los clicks
// de todos los botones "Eliminar", incluso los que se crean después.
function escucharBorrado(contenedor, url, recargarLista) {
    contenedor.addEventListener('click', async (evento) => {
        const boton = evento.target.closest('.btn-eliminar');
        if (!boton) return; // el click no fue en un botón de borrar

        const nombre = boton.closest('.tarjeta').querySelector('.tarjeta-titulo').textContent;
        if (!confirm(`¿Eliminar "${nombre}"?`)) return;

        boton.disabled = true;
        try {
            await eliminar(url, boton.dataset.id);
        } catch (error) {
            alert(`No se pudo eliminar: ${error.message}`);
        }
        await recargarLista();
    });
}

escucharBorrado(listaCueros, API_CUEROS, cargarCueros);
escucharBorrado(listaProductos, API_PRODUCTOS, cargarProductos);


// ---------------------------------------------------------------------
// Inicio: cargar las dos listas al abrir la página
// ---------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
    cargarCueros();
    cargarProductos();
});
