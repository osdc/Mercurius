/****  HEADER ****/
const hamburger = document.querySelector('.mobile-menu-button');
const menu = document.querySelector('.mobile-menu');

// Event listener
hamburger.addEventListener('click', () => {
  menu.classList.toggle('hidden');
});



