use std::sync::{Arc, Mutex};

use axum::{response::IntoResponse, routing, Router};

struct AppState {
    next_user_id: u32,
    lobbies: Vec<Lobby>,
}

#[tokio::main]
async fn main() {
    let state = AppState {
        next_user_id: 0,
        lobbies: Vec::new(),
    };

    let app = Router::new()
        .route("/script/:file", routing::get(routes::script))
        .route("/style/:file", routing::get(routes::style))
        .route("/", routing::get(routes::index))
        .with_state(Arc::new(Mutex::new(state)));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

mod html;

mod game;

mod routes;

enum Lobby {
    Open { users: Vec<String> },
    InGame,
}

struct LobbyUser {
    id: u32,
    name: String,
    role_preferences: Vec<String>,
}
