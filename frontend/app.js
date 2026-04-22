'use strict';
const API = 'http://localhost:8080';

const App = (() => {
  let token = localStorage.getItem('eh_token');
  let user  = JSON.parse(localStorage.getItem('eh_user') || 'null');

  // ── NAVIGATION ──────────────────────────────────────
  function showPage(name) {
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.nav-link').forEach(l => l.classList.remove('active'));
    const pg = document.getElementById('page-' + name);
    if (!pg) return;
    pg.classList.add('active');
    pg.classList.remove('page-enter'); void pg.offsetWidth; pg.classList.add('page-enter');
    document.querySelectorAll('.nav-link[data-page="' + name + '"]').forEach(l => l.classList.add('active'));
    if (name === 'events')  loadEvents();
    if (name === 'tickets') loadMyTickets();
    if (name === 'admin')   loadUsers();
  }
  document.querySelectorAll('.nav-link').forEach(l => {
    l.addEventListener('click', e => { e.preventDefault(); showPage(l.dataset.page); });
  });

  // ── TOAST ───────────────────────────────────────────
  function toast(msg, type = 'success') {
    const t = document.getElementById('toast');
    t.textContent = type === 'success' ? '✓  ' + msg : '✕  ' + msg;
    t.className = 'toast ' + type;
    clearTimeout(t._t);
    t._t = setTimeout(() => t.className = 'toast hidden', 3500);
  }

  // ── ALERT ───────────────────────────────────────────
  function showAlert(id, msg, type) {
    const el = document.getElementById(id);
    if (!el) return;
    el.textContent = msg; el.className = 'alert ' + type;
    clearTimeout(el._t); el._t = setTimeout(() => el.className = 'alert hidden', 4500);
  }

  // ── UPDATE NAV AFTER LOGIN ───────────────────────────
  function updateNav() {
    const authEl  = document.getElementById('navAuth');
    const userEl  = document.getElementById('navUserInfo');
    const adminLk = document.querySelector('.admin-link');
    const tickLk  = document.getElementById('ticketsLink');
    const createBtn = document.getElementById('createEventBtn');

    if (!user || !token) {
      authEl?.classList.remove('hidden');  userEl?.classList.add('hidden');
      adminLk?.style && (adminLk.style.display = 'none');
      tickLk?.style  && (tickLk.style.display  = 'none');
      if (createBtn) createBtn.classList.add('hidden');
      return;
    }
    authEl?.classList.add('hidden'); userEl?.classList.remove('hidden');
    document.getElementById('navAvatar').textContent  = user.first_name?.[0]?.toUpperCase() || '?';
    document.getElementById('navUserName').textContent = `${user.first_name} ${user.last_name}`;
    document.getElementById('navUserRole').textContent = user.role;
    if (adminLk) adminLk.style.display = user.role === 'admin' ? '' : 'none';
    if (tickLk)  tickLk.style.display  = 'block';
    if (createBtn) {
      createBtn.classList.toggle('hidden', user.role !== 'organizer' && user.role !== 'admin');
    }
  }

  // ── AUTH ────────────────────────────────────────────
  async function doLogin() {
    const btn   = document.getElementById('loginBtn');
    const email = document.getElementById('login-email').value;
    const pass  = document.getElementById('login-pass').value;
    if (!email || !pass) return showAlert('login-alert', 'Заполни все поля', 'error');
    btn.textContent = 'Загрузка...'; btn.disabled = true;
    try {
      const res  = await fetch(API + '/auth/login', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({ email, password: pass }) });
      const data = await res.json();
      btn.textContent = 'Войти'; btn.disabled = false;
      if (!res.ok) return showAlert('login-alert', data.error || 'Неверный email или пароль', 'error');
      token = data.token; user = data.user;
      localStorage.setItem('eh_token', token);
      localStorage.setItem('eh_user', JSON.stringify(user));
      updateNav();
      toast(`Привет, ${user.first_name}! 👋`);
      setTimeout(() => showPage('events'), 700);
    } catch {
      btn.textContent = 'Войти'; btn.disabled = false;
      showAlert('login-alert', '❌ Сервер недоступен. Запусти docker-compose up', 'error');
    }
  }

  async function doRegister() {
    const btn = document.getElementById('regBtn');
    const body = {
      first_name: document.getElementById('reg-fname').value,
      last_name:  document.getElementById('reg-lname').value,
      email:      document.getElementById('reg-email').value,
      password:   document.getElementById('reg-pass').value,
    };
    if (!body.first_name || !body.email || !body.password) return showAlert('reg-alert', 'Заполни все поля', 'error');
    btn.textContent = 'Загрузка...'; btn.disabled = true;
    try {
      const res  = await fetch(API + '/auth/register', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(body) });
      const data = await res.json();
      btn.textContent = 'Создать аккаунт'; btn.disabled = false;
      if (!res.ok) return showAlert('reg-alert', data.error || 'Ошибка регистрации', 'error');
      token = data.token; user = data.user;
      localStorage.setItem('eh_token', token);
      localStorage.setItem('eh_user', JSON.stringify(user));
      updateNav();
      toast('Аккаунт создан! 🎉');
      setTimeout(() => showPage('events'), 700);
    } catch {
      btn.textContent = 'Создать аккаунт'; btn.disabled = false;
      showAlert('reg-alert', '❌ Сервер недоступен', 'error');
    }
  }

  function doLogout() {
    token = null; user = null;
    localStorage.removeItem('eh_token'); localStorage.removeItem('eh_user');
    updateNav(); toast('Вы вышли', 'error'); showPage('home');
  }

  // ── EVENTS ──────────────────────────────────────────
  async function loadEvents() {
    const grid = document.getElementById('events-grid');
    const countEl = document.getElementById('eventsCount');
    if (!token) {
      grid.innerHTML = `<div class="empty-state"><div class="empty-icon">🔐</div><p>Войдите чтобы видеть события</p><button class="btn-primary" onclick="App.showPage('login')">Войти</button></div>`;
      return;
    }
    grid.innerHTML = `<div class="loading-state">Загрузка событий...</div>`;
    try {
      const res  = await fetch(API + '/events', { headers:{ Authorization:'Bearer ' + token } });
      const data = await res.json();
      if (!res.ok || !Array.isArray(data) || data.length === 0) {
        countEl.textContent = 'Нет событий';
        grid.innerHTML = `<div class="empty-state"><div class="empty-icon">📅</div><p>Событий пока нет</p>${user?.role==='organizer'||user?.role==='admin'?'<button class="btn-primary" onclick="App.toggleCreateForm()">Создать первое событие</button>':''}</div>`;
        return;
      }
      countEl.textContent = `${data.length} событий`;
      grid.innerHTML = data.map((e, i) => `
        <div class="event-card" style="animation-delay:${i*.06}s">
          <div class="event-card-header">
            <span class="event-num">EVENT #${e.id}</span>
            <div class="event-name">${e.title}</div>
            <div class="event-info">
              <span>📍 Площадка #${e.venue_id}</span>
              <span>🕐 ${new Date(e.start_time).toLocaleString('ru',{day:'numeric',month:'long',hour:'2-digit',minute:'2-digit'})}</span>
            </div>
          </div>
          <div class="event-desc-text">${e.description || 'Описание не указано'}</div>
          <div class="event-card-footer">
            <button class="btn-buy" onclick="App.bookTicket(${e.id}, '${e.title.replace(/'/g,"\\'")}')">🎫 Купить билет</button>
          </div>
        </div>
      `).join('');
    } catch {
      grid.innerHTML = `<div class="empty-state"><div class="empty-icon">❌</div><p>Сервер недоступен</p></div>`;
    }
  }

  function toggleCreateForm() {
    document.getElementById('createFormBox').classList.toggle('hidden');
  }

  async function doCreateEvent() {
    const body = {
      title:       document.getElementById('ev-title').value,
      description: document.getElementById('ev-desc').value,
      venue_id:    parseInt(document.getElementById('ev-venue').value),
      start_time:  new Date(document.getElementById('ev-time').value).toISOString(),
    };
    if (!body.title || !body.start_time) return showAlert('create-event-alert', 'Заполни название и дату', 'error');
    try {
      const res  = await fetch(API + '/events', { method:'POST', headers:{'Content-Type':'application/json', Authorization:'Bearer ' + token}, body: JSON.stringify(body) });
      const data = await res.json();
      if (!res.ok) return showAlert('create-event-alert', data.error || 'Ошибка', 'error');
      toggleCreateForm(); toast('Событие создано! 🎉'); loadEvents();
    } catch {
      showAlert('create-event-alert', '❌ Ошибка', 'error');
    }
  }

  // ── TICKETS ─────────────────────────────────────────
  async function bookTicket(eventId, eventTitle) {
    if (!token) { toast('Войди чтобы купить билет', 'error'); showPage('login'); return; }
    try {
      const res  = await fetch(API + '/tickets/book', { method:'POST', headers:{'Content-Type':'application/json', Authorization:'Bearer ' + token}, body: JSON.stringify({ event_id: eventId }) });
      const data = await res.json();
      if (!res.ok) return toast(data.error || 'Ошибка бронирования', 'error');
      toast(`Билет куплен! 🎫`);
      if (data.qr_code) showQR(data, eventTitle);
      else showPage('tickets');
    } catch {
      toast('Сервер недоступен', 'error');
    }
  }

  async function loadMyTickets() {
    const cont = document.getElementById('tickets-container');
    if (!token) {
      cont.innerHTML = `<div class="empty-state"><div class="empty-icon">🔐</div><p>Войдите чтобы видеть билеты</p><button class="btn-primary" onclick="App.showPage('login')">Войти</button></div>`;
      return;
    }
    cont.innerHTML = `<div class="loading-state">Загрузка билетов...</div>`;
    try {
      const res  = await fetch(API + '/tickets/my', { headers:{ Authorization:'Bearer ' + token } });
      const data = await res.json();
      if (!res.ok || !Array.isArray(data) || data.length === 0) {
        cont.innerHTML = `<div class="empty-state"><div class="empty-icon">🎫</div><p>Нет билетов</p><button class="btn-primary" onclick="App.showPage('events')">Найти событие</button></div>`;
        return;
      }
      cont.innerHTML = `<div class="tickets-grid">${data.map((t, i) => `
        <div class="ticket-card" style="animation-delay:${i*.06}s">
          <div class="ticket-card-body">
            <div class="ticket-header">
              <span class="ticket-num">TICKET #${t.id}</span>
              <span class="ticket-status">✓ Активен</span>
            </div>
            <div class="ticket-event-name">${t.event_title || 'Событие #' + t.event_id}</div>
            <div class="ticket-date">📅 ${new Date(t.created_at).toLocaleString('ru',{day:'numeric',month:'long',year:'numeric',hour:'2-digit',minute:'2-digit'})}</div>
          </div>
          <div class="ticket-card-footer">
            ${t.qr_code ? `<button class="btn-show-qr" onclick="App.showQR(${JSON.stringify(t).replace(/"/g,'&quot;')}, '${(t.event_title||'Событие').replace(/'/g,"\\'")}')">📱 Показать QR-код</button>` : '<span style="color:var(--muted);font-size:.82rem">QR недоступен</span>'}
            <span class="ticket-id-small">Event #${t.event_id}</span>
          </div>
        </div>
      `).join('')}</div>`;
    } catch {
      cont.innerHTML = `<div class="empty-state"><div class="empty-icon">❌</div><p>Сервер недоступен</p></div>`;
    }
  }

  // ── QR MODAL ────────────────────────────────────────
  function showQR(ticket, eventTitle) {
    document.getElementById('qrTitle').textContent = eventTitle || 'Твой билет';
    document.getElementById('qrInfo').textContent  = `Билет #${ticket.id} · Событие #${ticket.event_id}`;
    document.getElementById('qrImg').src = ticket.qr_code || '';
    document.getElementById('qrModal').classList.remove('hidden');
  }
  function closeQR() { document.getElementById('qrModal').classList.add('hidden'); }

  // ── ADMIN ───────────────────────────────────────────
  async function loadUsers() {
    const cont = document.getElementById('admin-container');
    const cnt  = document.getElementById('usersCount');
    if (!token) { cont.innerHTML = `<div class="empty-state"><div class="empty-icon">🔐</div><p>Нет доступа</p></div>`; return; }
    try {
      const res  = await fetch(API + '/admin/users', { headers:{ Authorization:'Bearer ' + token } });
      const data = await res.json();
      if (!res.ok) { cont.innerHTML = `<div class="empty-state"><div class="empty-icon">❌</div><p>${data.error}</p></div>`; return; }
      cnt.textContent = `${data.length} пользователей`;
      cont.innerHTML = `
        <table class="admin-table">
          <thead><tr><th>ID</th><th>Имя</th><th>Email</th><th>Роль</th><th>Дата регистрации</th></tr></thead>
          <tbody>${data.map(u => `
            <tr>
              <td style="font-family:monospace;color:var(--muted)">#${u.id}</td>
              <td><strong>${u.first_name} ${u.last_name}</strong></td>
              <td style="color:var(--muted)">${u.email}</td>
              <td><span class="role-tag ${u.role}">${u.role}</span></td>
              <td style="color:var(--muted)">${new Date(u.created_at).toLocaleDateString('ru',{day:'numeric',month:'long',year:'numeric'})}</td>
            </tr>
          `).join('')}</tbody>
        </table>`;
    } catch {
      cont.innerHTML = `<div class="empty-state"><div class="empty-icon">❌</div><p>Сервер недоступен</p></div>`;
    }
  }

  // ── INIT ────────────────────────────────────────────
  updateNav();
  if (token) {
    fetch(API + '/me', { headers:{ Authorization:'Bearer ' + token } })
      .then(r => r.ok ? r.json() : null).then(d => { if (d) { user = d; localStorage.setItem('eh_user', JSON.stringify(d)); updateNav(); } })
      .catch(() => {});
  }

  return { showPage, doLogin, doRegister, doLogout, loadEvents, loadMyTickets, loadUsers, bookTicket, toggleCreateForm, doCreateEvent, showQR, closeQR };
})();