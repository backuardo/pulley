const PAGE_SIZE = 20;
let currentPaginationOffset = 0;

const Controller = {
  search: async (ev) => {
    ev.preventDefault();
    const query = document.getElementById('query').value;
    try {
      await Controller.fetchAndDisplayResults(query, 0);
    } catch (error) {
      Controller.showError(error.message);
    }
  },

  showError: (message) => {
    const errorElement = document.getElementById('error-message');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
  },

  updateTable: (results) => {
    const table = document.getElementById('table-body');
    table.innerHTML = '';

    for (let result of results) {
      const row = document.createElement('tr');
      const cell = document.createElement('td');
      cell.textContent = result;
      row.appendChild(cell);
      table.appendChild(row);
    }
  },

  appendToTable: (results) => {
    const table = document.getElementById('table-body');
    for (let result of results) {
      const row = document.createElement('tr');
      const cell = document.createElement('td');
      cell.textContent = result;
      row.appendChild(cell);
      table.appendChild(row);
    }
  },

  loadMore: async () => {
    const query = document.getElementById('query').value;
    currentPaginationOffset += PAGE_SIZE;
    await Controller.fetchAndDisplayResults(query, currentPaginationOffset);
  },

  fetchAndDisplayResults: async (query, offset) => {
    try {
      const response = await fetch(
        `/search?q=${encodeURIComponent(
          query
        )}&offset=${offset}&limit=${PAGE_SIZE}`
      );
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const results = await response.json();
      offset === 0
        ? Controller.updateTable(results)
        : Controller.appendToTable(results);
    } catch (error) {
      Controller.showError(error.message);
    }
  },
};

document.getElementById('form').addEventListener('submit', Controller.search);
document
  .getElementById('load-more')
  .addEventListener('click', Controller.loadMore);
