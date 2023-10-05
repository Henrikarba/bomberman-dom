import { FrameWork } from "./main"

class Router {
    constructor(container){
        this.routes = {}
        this.container = container
        window.addEventListener('popstate', this.handleRouteChange.bind(this));
    }

    registerRoute(path, handler) {
        this.routes[path] = handler;
      }
    
    handleRouteChange() {
        const path = window.location.pathname;
        const element = this.routes[path];
    
        if (element) {  
            FrameWork.InitFramework(element(), this.container)
        } else {
            const errorMessage = document.createElement('h1');
            errorMessage.textContent = '404 - Page Not Found';
            this.container.appendChild(errorMessage);
        }
    }
}

export default Router
