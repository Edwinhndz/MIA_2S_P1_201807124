import React from "react";
import { Link } from "react-router-dom";

function NavBar() {
    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-#460b0b">
            <div className="container-fluid">
                <Link to="/" className="navbar-brand"></Link>
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                    <span className="navbar-toggler-icon"></span>
                </button>
              
            </div>
        </nav>
    );
}

export default NavBar;
