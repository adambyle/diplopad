use std::{fs::File, io::Read};

use axum::{
    extract,
    http::{header, StatusCode},
    response::{IntoResponse, Response},
};

pub async fn index() -> Response {
    "".into_response()
}

pub async fn script(extract::Path(file): extract::Path<String>) -> Response {
    static_file(&format!("./scripts/{file}"), "application/javascript")
}

pub async fn style(extract::Path(file): extract::Path<String>) -> Response {
    static_file(&format!("./styles/{file}"), "text/css")
}

fn static_file(path: &str, content_type: &str) -> Response {
    let mut contents = Vec::new();
    let mut file = match File::open(path) {
        Ok(file) => file,
        Err(_) => return StatusCode::NOT_FOUND.into_response(),
    };
    file.read_to_end(&mut contents).unwrap();
    ([(header::CONTENT_TYPE, content_type)], contents).into_response()
}
