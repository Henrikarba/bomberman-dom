import Router from './framework/router';
import { FrameWork } from './framework/main';

const container = document.getElementById('app-container');
const router = new Router(container);

// Register your routes
router.registerRoute('/', () => FrameWork.CreateElement('h1', null, 'Home Page'));
router.registerRoute('/about', () => FrameWork.CreateElement('h1', null, 'About Page'));
// Add more routes as needed

const homeButton = document.getElementById('home-button'); // Replace with your element ID
console.log(homeButton)
homeButton.addEventListener('click', () => {
    window.history.pushState(null, '', '/about');
    router.handleRouteChange();
});