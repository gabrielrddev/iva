<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';
  import { goto } from '$app/navigation';
  const dispatch = createEventDispatcher();


  //TODO tenho que pegar para quando for criado colocar o email no localstorage e fazer o get automatico
  let showSignup = true;
  let name = '';
  let email = '';
  let password = '';
  let repeatPassword = '';
  let phone = '';
  let numerQuestion = '';
  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
  const CREATE_USER_URL = `${API_BASE_URL}/users/create`;

  let notificationMessage = '';
  let notificationType = 'info';
  let notificationVisible = false;
  let notificationTimer;

  const notify = (message, type = 'info') => {
    notificationMessage = message;
    notificationType = type;
    notificationVisible = true;
    clearTimeout(notificationTimer);
    notificationTimer = setTimeout(() => {
      notificationVisible = false;
    }, 3500);
  };

  // Quando o popup abrir, bloqueia o scroll
  onMount(() => {
    if (browser) {
      document.body.style.overflow = 'hidden';
    }
  });

  // Quando o popup fechar, desbloqueia o scroll
  const closePopup = () => {
    if (browser) {
      document.body.style.overflow = '';
    }
    dispatch('close');
  };

  const closePopupAndGoToHome = () => {
    closePopup();
    goto('/');
  };

  const toggleForm = () => {
    showSignup = !showSignup;
    name = '';
    email = '';
    phone = '';
    password = '';
    repeatPassword = '';
    numerQuestion = '0';
  };

  const handleSubmit = async () => {
    if (!email || !password || (showSignup && !name)) {
      notify('Preencha todos os campos!', 'warning');
      return;
    }

    if (showSignup && password !== repeatPassword) {
      notify('As senhas não conferem!', 'error');
      return;
    }

    try {
      if (showSignup) {
        const res = await fetch(CREATE_USER_URL, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name, email, phone, password, numerQuestion })
        });
        const data = await res.json().catch(() => null);

        if (res.ok && data && data.status === 'success') {
          if (browser) {
            localStorage.setItem('cnt_email', email);
          }
          notify('Conta criada com sucesso!', 'success');
          if (browser) {
            window.location.reload();
          }
          closePopup();
        } else {
          const msg = (data && data.message) ? data.message : 'Erro ao criar usuário';
          notify(msg, 'error');
        }
      } else {
        const res = await fetch(`${API_BASE_URL}/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password })
        });
        const data = await res.json().catch(() => null);
        if (res.ok && data && data.status === 'success') {
          if (browser) {
            localStorage.setItem('cnt_email', email);
            window.location.reload();
          }
          notify('Login realizado com sucesso!', 'success');
          closePopup();
        } else if (data && data.status === 'user_not_found') {
          notify('Usuário não encontrado', 'error');
        } else if (data && data.status === 'wrong_password') {
          notify('Senha incorreta', 'error');
        } else {
          notify('Erro ao realizar login', 'error');
        }
      }
    } catch (err) {
      notify('Erro ao criar usuário', 'error');
    }
  };

  // Garantir que, se o componente for destruído, o scroll volte
  onDestroy(() => {
    if (browser) {
      document.body.style.overflow = '';
    }
  });
</script>

<!-- POPUP FULLSCREEN COM FUNDO Borrado -->
{#if notificationVisible}
  <div class="fixed top-4 right-4 z-50">
    <div class="px-4 py-3 rounded-lg shadow-md text-sm text-white"
      class:bg-[#00995D]={notificationType === 'success'}
      class:bg-red-600={notificationType === 'error'}
      class:bg-yellow-500={notificationType === 'warning'}
      class:bg-gray-800={notificationType === 'info'}>
      {notificationMessage}
    </div>
  </div>
{/if}
<div class="fixed inset-0 flex items-center justify-center z-50 mt-20">
  
    <!-- Fundo borrado -->
  <button
    class="absolute inset-0  bg-opacity-30 backdrop-blur-sm"
    aria-label="Fechar popup"
    on:click={closePopupAndGoToHome}>
  </button>

  <!-- Modal -->
  <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md p-8 z-10">
    <button
      class="absolute top-4 right-4 text-gray-500 hover:text-gray-800 text-lg font-bold"
      on:click={closePopupAndGoToHome}>
      &times;
    </button>

    <h2 class="text-2xl font-semibold text-gray-800 mb-6 text-center">
      {showSignup ? 'Criar Conta' : 'Login'}
    </h2>

    <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
      {#if showSignup}
      <div>
        <label for="name" class="block text-sm font-medium text-gray-700">Nome</label>
        <input type="text" bind:value={name} placeholder="Seu nome"
          class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"/>
      </div>
      {/if}

      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input type="email" bind:value={email} placeholder="seu@email.com"
          class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"/>
      </div>
      {#if showSignup}
        <div>
        <label for="phone" class="block text-sm font-medium text-gray-700">Phone</label>
        <input type="phone" bind:value={phone} placeholder="(48) 99999-9999"
          class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"/>
      </div>
{/if}
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Senha</label>
        <input type="password" bind:value={password} placeholder="Senha"
          class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"/>
      </div>

      {#if showSignup}
      <div>
        <label for="repeatPassword" class="block text-sm font-medium text-gray-700">Repetir Senha</label>
        <input type="password" bind:value={repeatPassword} placeholder="Repita a senha"
          class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"/>
      </div>
      {/if}

      <button type="submit"
        class="cursor-pointer w-full bg-[#00995D] text-white py-2 rounded-lg text-sm hover:bg-[#00814F] transition mt-4">
        {showSignup ? 'Criar Conta' : 'Entrar'}
      </button>
    </form>

    <p class="text-center text-sm text-gray-600 mt-4">
      {showSignup ? 'Já tem conta?' : 'Não tem conta?'}
      <button class="cursor-pointer text-[#00995D] font-semibold hover:underline ml-1" on:click={toggleForm}>
        {showSignup ? 'Login' : 'Criar Conta'}
      </button>
    </p>
  </div>
</div>
