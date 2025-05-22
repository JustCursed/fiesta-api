document.addEventListener('DOMContentLoaded', () => {
  // Текущие выбранные даты
  let selectedDates = [];

  // Функция для получения параметров URL
  function getUrlParams() {
    const params = new URLSearchParams(window.location.search);
    const datesParam = params.get('dates');
    if (datesParam) {
      if (datesParam.includes('--')) {
        return {
          dates: datesParam.split('--').sort((a, b) => new Date(b) - new Date(a))
        };
      }
      return {
        dates: [datesParam]
      };
    }
    return { dates: [] };
  }

  // Функция для обновления URL
  function updateUrl() {
    let newUrl = window.location.pathname;
    if (selectedDates.length === 1) {
      newUrl += `?dates=${selectedDates[0]}`;
    } else if (selectedDates.length === 2) {
      // Сортируем даты - сначала более свежая
      const sortedDates = [...selectedDates].sort((a, b) => new Date(b) - new Date(a));
      newUrl += `?dates=${sortedDates[0]}--${sortedDates[1]}`;
    } else {
      newUrl = window.location.pathname;
    }
    window.history.pushState({}, '', newUrl);
  }

  // Переключение боковой панели
  const categoryIcon = document.querySelector('.category-icon');
  const sidebar = document.querySelector('.sidebar');

  if (categoryIcon && sidebar) {
    categoryIcon.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      sidebar.classList.toggle('collapsed');
    });
  }

  // Раскрытие/скрытие списков в боковой панели
  const toggles = document.querySelectorAll('.nav-section-toggle');
  toggles.forEach(toggle => {
    toggle.addEventListener('click', (event) => {
      event.preventDefault();
      const collapsibleItem = toggle.closest('.collapsible-item');
      if (!collapsibleItem) return;

      const collapsibleList = collapsibleItem.querySelector('.collapsible-list');
      if (!collapsibleList) return;

      const toggleIcon = toggle.querySelector('.toggle-icon');
      collapsibleList.classList.toggle('expanded');
      if (toggleIcon) {
        toggleIcon.classList.toggle('rotated');
      }
    });
  });

  // Выбор сервера
  const serverItems = document.querySelectorAll('.collapsible-list li a:not(.date-item)');
  serverItems.forEach(item => {
    item.addEventListener('click', (event) => {
      event.preventDefault();

      const parentList = item.closest('.collapsible-list');
      if (parentList) {
        parentList.querySelectorAll('li a').forEach(i => i.classList.remove('active'));
      }
      item.classList.add('active');
    });
  });

  // Выбор даты
  const dateItems = document.querySelectorAll('.date-item');
  dateItems.forEach(item => {
    item.addEventListener('click', (event) => {
      event.preventDefault();
      const date = item.textContent.trim();
      const isCtrlPressed = event.ctrlKey || event.metaKey; // Cmd на Mac

      if (!isCtrlPressed) {
        // Обычный клик - сбрасываем все выделения
        selectedDates = [date];
        dateItems.forEach(i => i.classList.remove('active'));
        item.classList.add('active');
      } else {
        // Клик с Ctrl/Cmd - добавляем/удаляем дату
        if (selectedDates.includes(date)) {
          // Удаляем дату из выбранных
          selectedDates = selectedDates.filter(d => d !== date);
          item.classList.remove('active');
        } else {
          // Добавляем дату, но не более 2
          if (selectedDates.length < 2) {
            selectedDates.push(date);
            item.classList.add('active');
          } else {
            // Если уже 2 даты выбрано, заменяем последнюю
            const lastSelected = selectedDates[1] || selectedDates[0];
            const lastSelectedItem = [...dateItems].find(i => i.textContent.trim() === lastSelected);
            if (lastSelectedItem) {
              lastSelectedItem.classList.remove('active');
            }
            selectedDates = [selectedDates[0], date];
            item.classList.add('active');
          }
        }
      }

      // Обновляем URL и логи
      updateUrl();
      updateLogsWithDates();
    });
  });

  // Настройки (выпадающее меню)
  const settingsButton = document.querySelector('.settings-button');
  if (settingsButton) {
    settingsButton.addEventListener('click', (e) => {
      e.stopPropagation();
      settingsButton.classList.toggle('active');
    });

    document.addEventListener('click', (e) => {
      const settingsDropdown = settingsButton.querySelector('.settings-dropdown');
      if (settingsButton.classList.contains('active') && settingsDropdown && !settingsButton.contains(e.target) && !settingsDropdown.contains(e.target)) {
        settingsButton.classList.remove('active');
      }
    });

    const settingsDropdown = settingsButton.querySelector('.settings-dropdown');
    if (settingsDropdown) {
      settingsDropdown.addEventListener('click', (e) => {
        e.stopPropagation();
      });
    }
  }

  // Загрузка логов
  const logsListElement = document.getElementById('logs-list');

  if (logsListElement) {
    // Функция для обновления URL запроса с учетом дат
    function getLogsUrl() {
      const baseUrl = '/api/logs/htc/item';
      if (selectedDates.length === 0) {
        return `${baseUrl}?data=2021-03-10--2023-03-10&params=player&message=wdawdawdaw`;
      } else if (selectedDates.length === 1) {
        return `${baseUrl}?data=${selectedDates[0]}--${selectedDates[0]}&params=player&message=wdawdawdaw`;
      } else {
        // Сортируем даты - сначала более свежая
        const sortedDates = [...selectedDates].sort((a, b) => new Date(b) - new Date(a));
        return `${baseUrl}?data=${sortedDates[0]}--${sortedDates[1]}&params=player&message=wdawdawdaw`;
      }
    }

    // Функция для обновления логов с учетом выбранных дат
    function updateLogsWithDates() {
      const fetchUrl = getLogsUrl();
      fetchAndDisplayLogs(fetchUrl);
    }

    async function fetchAndDisplayLogs(fetchUrl) {
      try {
        logsListElement.innerHTML = '<div class="message">Загрузка логов...</div>';

        const response = await fetch(fetchUrl);

        if (!response.ok) {
          const errorText = `Ошибка HTTP: статус ${response.status}`;
          logsListElement.innerHTML = `<div class="message">Не удалось загрузить логи: ${errorText}</div>`;
          return;
        }

        const data = await response.json();

        logsListElement.innerHTML = '';

        if (data && Array.isArray(data.logs) && data.logs.length > 0) {
          data.logs.forEach(logItem => {
            const messageElement = document.createElement('div');
            messageElement.classList.add('message');
            messageElement.textContent = String(logItem).trim();
            logsListElement.appendChild(messageElement);
          });
        } else {
          const messageElement = document.createElement('div');
          messageElement.classList.add('message');
          messageElement.textContent = 'Логи не найдены.';
          logsListElement.appendChild(messageElement);
        }
      } catch (error) {
        logsListElement.innerHTML = `<div class="message">Не удалось загрузить логи: ${error.message}</div>`;
      }
    }

    // При загрузке страницы проверяем параметры URL
    const urlParams = getUrlParams();
    selectedDates = urlParams.dates || [];

    // Отмечаем выбранные даты в интерфейсе
    if (selectedDates.length > 0) {
      dateItems.forEach(item => {
        const date = item.textContent.trim();
        if (selectedDates.includes(date)) {
          item.classList.add('active');
        }
      });

      // Раскрываем список дат, если он свернут
      const firstDateItem = [...dateItems].find(item => selectedDates.includes(item.textContent.trim()));
      if (firstDateItem) {
        const collapsibleItem = firstDateItem.closest('.collapsible-item');
        if (collapsibleItem) {
          const collapsibleList = collapsibleItem.querySelector('.collapsible-list');
          const toggle = collapsibleItem.querySelector('.nav-section-toggle');
          const toggleIcon = toggle?.querySelector('.toggle-icon');

          if (collapsibleList) collapsibleList.classList.add('expanded');
          if (toggleIcon) toggleIcon.classList.add('rotated');
        }
      }
    }

    // Загружаем логи
    updateLogsWithDates();
  }

  // Обработка изменения истории (назад/вперед в браузере)
  window.addEventListener('popstate', () => {
    const urlParams = getUrlParams();
    selectedDates = urlParams.dates || [];

    // Обновляем активные даты в интерфейсе
    dateItems.forEach(item => {
      const date = item.textContent.trim();
      item.classList.toggle('active', selectedDates.includes(date));
    });

    // Обновляем логи
    updateLogsWithDates();
  });
});
