<script>
  import Header from '$lib/components/layout/Header.svelte';
  import whatsapp from '$lib/assets/whatsapp.png';
  import PopUp from '$lib/components/layout/PopUp.svelte';
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

  let oldPassword = '';
  let newPassword = '';

  let keepConnect = false;
  if (browser) {
    keepConnect = !!localStorage.getItem('cnt_email');
  }

  let userName = '';
  let userEmail = '';
  let userPhone = '';
  let originalEmail = '';

  let isEditing = false;

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

  onMount(async () => {
    if (!browser) return;
    const email = localStorage.getItem('cnt_email');
    if (!email) return;
    try {
      const res = await fetch(`${API_BASE_URL}/users/get?email=${encodeURIComponent(email)}`);
      const data = await res.json().catch(() => null);
      if (res.ok && data && !data.status) {
        userName = data.name || '';
        userEmail = data.email || '';
        userPhone = data.phone || '';
        originalEmail = data.email || email;
      }
    } catch (e) {
      // swallow fetch errors
    }
  });
  const handlePasswordChange = async () => {
    if (!oldPassword || !newPassword) {
      notify('Preencha todos os campos.', 'warning');
      return;
    }
    if (!browser) return;
    const email = originalEmail || userEmail || localStorage.getItem('cnt_email') || '';
    if (!email) {
      notify('Email não encontrado. Faça login novamente.', 'error');
      return;
    }
    try {
      const res = await fetch(`${API_BASE_URL}/users/password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, oldPassword, newPassword })
      });
      const data = await res.json().catch(() => null);
      console.log(data)
      if (res.ok && data && data.status === 'success') {
        notify('Senha alterada com sucesso!', 'success');
        oldPassword = '';
        newPassword = '';
      } else if (data && data.status === 'user_not_found') {
        notify('Usuário não encontrado', 'error');
      } else if (data && data.status === 'wrong_password') {
        notify('Senha atual incorreta', 'error');
      } else {
        notify('Erro ao alterar senha', 'error');
      }
    } catch (e) {
      notify('Erro ao alterar senha', 'error');
    }
  };

  const handleEditOrSubmit = async () => {
    if (!isEditing) {
      isEditing = true;
      return;
    }
    if (!browser) return;
    const payload = {
      email: originalEmail || userEmail,
      newName: userName,
      newEmail: userEmail,
      newPhone: userPhone
    };
    try {
      const res = await fetch(`${API_BASE_URL}/users/edit`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      const data = await res.json().catch(() => null);
      if (res.ok && data && data.status === 'success') {
        notify('Atualizado com sucesso!', 'success');
        // se email mudou, atualiza localStorage para manter sessão
        const current = localStorage.getItem('cnt_email');
        if (payload.newEmail && payload.newEmail !== current) {
          localStorage.setItem('cnt_email', payload.newEmail);
        }
        isEditing = false;
      } else {
        const msg = (data && data.message) ? data.message : 'Erro ao atualizar';
        notify(msg, 'error');
      }
    } catch (e) {
      notify('Erro ao atualizar', 'error');
    }
  };

  const openWhatsAppChat = () => {
    const message = encodeURIComponent('Olá! Iva, preciso de ajuda.');
    const whatsappLink = `https://wa.me/5548984684188?text=${message}`;
    window.open(whatsappLink, '_blank');
  };

  const handleLogout = () => {
    if (!browser) return;
    localStorage.removeItem('cnt_email');
    window.location.href = '/';
  };
  let questionText = '';
  let chatAnswer = '';
  let chatLoading = false;

  const handleAskAssistant = async () => {
    if (!browser) return;
    const email = originalEmail || userEmail || localStorage.getItem('cnt_email') || '';
    const question = questionText.trim();
    if (!email || !question) {
      alert('Preencha a pergunta.');
      return;
    }
    chatLoading = true;
    chatAnswer = '';
    try {
      const res = await fetch(`${API_BASE_URL}/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, message: question })
      });
      const data = await res.json().catch(() => null);
      let answer = '';
      if (data) {
        if (data.status === 'success' && data.answer) answer = data.answer;
        else if (data.response) answer = data.response;
        else if (data.answer) answer = data.answer;
      }
      if (res.ok && answer) {
        chatAnswer = answer;
      } else {
        const msg = (data && (data.message || data.error)) ? (data.message || data.error) : 'Erro ao enviar pergunta';
        notify(msg, 'error');
      }
    } catch (e) {
      notify('Erro ao enviar pergunta', 'error');
    } finally {
      chatLoading = false;
    }
  };
</script>

<!-- HEADER -->

{#if keepConnect}
<Header />

<!-- Notification -->
{#if notificationVisible}
  <div class="fixed top-20 right-4 z-50">
    <div class="px-4 py-3 rounded-lg shadow-md text-sm text-white"
      class:bg-[#00995D]={notificationType === 'success'}
      class:bg-red-600={notificationType === 'error'}
      class:bg-yellow-500={notificationType === 'warning'}
      class:bg-gray-800={notificationType === 'info'}>
      {notificationMessage}
    </div>
  </div>
{/if}

<!-- CONTEÚDO PRINCIPAL -->
<main class="pt-24 min-h-screen bg-gray-50 flex justify-center px-4 pb-16">
  <div class="w-full max-w-5xl space-y-10">

    <!-- TÍTULO GERAL -->
    <div class="text-center mb-4">
      <h1 class="text-3xl font-semibold text-gray-800">Configurações da Conta</h1>
      <p class="text-gray-500 text-sm mt-2">Gerencie suas informações, segurança e suporte.</p>
    </div>

    <!-- SEÇÃO DADOS DO USUÁRIO -->
    <section class="bg-white shadow-lg rounded-2xl border border-gray-100 p-8">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Dados do Usuário</h2>
      <p class="text-gray-600 mb-6">Visualize e edite as informações associadas à sua conta.</p>

      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-700">Nome Completo</label>
          <input
            type="text"
            bind:value={userName}
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none bg-gray-50"
            disabled={!isEditing}
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
          <input
            type="email"
            bind:value={userEmail}
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none bg-gray-50"
            disabled={!isEditing}
          />
        </div>

        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700">Telefone</label>
          <input
            type="text"
            bind:value={userPhone}
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none bg-gray-50"
            disabled={!isEditing}
          />
        </div>
      </div>

      <div class="mt-8 text-right space-x-3">
        <button
          class="cursor-pointer bg-red-600 text-white py-2 px-6 rounded-lg text-sm hover:bg-red-700 transition"
          on:click={handleLogout}>
          Sair
        </button>
        <button
          class="cursor-pointer bg-[#00995D] text-white py-2 px-6 rounded-lg text-sm hover:bg-[#00814F] transition"
          on:click={handleEditOrSubmit}>
          {isEditing ? 'Realizar mudança' : 'Editar informações'}
        </button>
      </div>
    </section>

    <!-- SEÇÃO ALTERAR SENHA -->
    <section class="bg-white shadow-lg rounded-2xl border border-gray-100 p-8">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Alterar Senha</h2>
      <p class="text-gray-600 mb-6">
        Mantenha sua conta segura alterando sua senha regularmente.
      </p>

      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <label for="oldPassword" class="block text-sm font-medium text-gray-700">Senha Atual</label>
          <input
            type="password"
            bind:value={oldPassword}
            placeholder="Senha atual"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"
          />
        </div>

        <div>
          <label for="newPassword" class="block text-sm font-medium text-gray-700">Nova Senha</label>
          <input
            type="password"
            bind:value={newPassword}
            placeholder="Nova senha"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"
          />
        </div>
      </div>

      <div class="mt-8 text-right">
        <button
          class="cursor-pointer bg-[#00995D] text-white py-2 px-6 rounded-lg text-sm hover:bg-[#00814F] transition"
          on:click={handlePasswordChange}>
          Atualizar Senha
        </button>
      </div>
    </section>

    <!-- SEÇÃO SUPORTE (IVA WHATSAPP) -->
    <section class="shadow-lg rounded-2xl border border-gray-100 p-8 text-center" style="background-color: #E6F7EE;">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Iva - Suporte via WhatsApp</h2>
      <p class="text-gray-600 mb-6">
        Converse com nossa assistente virtual no WhatsApp e receba ajuda imediata.
      </p>

      <img
        src={whatsapp}
        alt="WhatsApp"
        class="w-20 mx-auto mb-4"
      />

      <button
        class="cursor-pointer bg-[#00995D] text-white py-3 px-8 rounded-lg text-sm hover:bg-[#00814F] transition"
        on:click={openWhatsAppChat}>
        Abrir Chat no WhatsApp
      </button>
    </section>
  </div>
</main>
{:else}
<PopUp />
<Header />
<!-- CONTEÚDO PRINCIPAL -->
<main class="pt-24 min-h-screen bg-gray-50 flex justify-center px-4 pb-16">
  <div class="w-full max-w-5xl space-y-10">

    <!-- TÍTULO GERAL -->
    <div class="text-center mb-4">
      <h1 class="text-3xl font-semibold text-gray-800">Configurações da Conta</h1>
      <p class="text-gray-500 text-sm mt-2">Gerencie suas informações, segurança e suporte.</p>
    </div>

    <!-- SEÇÃO DADOS DO USUÁRIO -->
    <section class="bg-white shadow-lg rounded-2xl border border-gray-100 p-8">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Dados do Usuário</h2>
      <p class="text-gray-600 mb-6">Visualize e edite as informações associadas à sua conta.</p>

      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-700">Nome Completo</label>
          <input
            type="text"
            placeholder="Gabriel Rodrigues Dias"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none bg-gray-50"
            disabled
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
          <input
            type="email"
            placeholder="gabriel@email.com"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none bg-gray-50"
            disabled
          />
        </div>

        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700">Telefone</label>
          <input
            type="text"
            placeholder="(48) 99999-9999"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none bg-gray-50"
            disabled
          />
        </div>
      </div>

      <div class="mt-8 text-right">
        <button
          class="cursor-pointer bg-[#00995D] text-white py-2 px-6 rounded-lg text-sm hover:bg-[#00814F] transition">
          Editar informações
        </button>
      </div>
    </section>

    <!-- SEÇÃO ALTERAR SENHA -->
    <section class="bg-white shadow-lg rounded-2xl border border-gray-100 p-8">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Alterar Senha</h2>
      <p class="text-gray-600 mb-6">
        Mantenha sua conta segura alterando sua senha regularmente.
      </p>

      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <label for="oldPassword" class="block text-sm font-medium text-gray-700">Senha Atual</label>
          <input
            type="password"
            bind:value={oldPassword}
            placeholder="Senha atual"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"
          />
        </div>

        <div>
          <label for="newPassword" class="block text-sm font-medium text-gray-700">Nova Senha</label>
          <input
            type="password"
            bind:value={newPassword}
            placeholder="Nova senha"
            class="w-full border rounded-lg px-3 py-2 mt-1 text-sm focus:ring-2 focus:ring-[#00995D] outline-none"
          />
        </div>
      </div>

      <div class="mt-8 text-right">
        <button
          class="cursor-pointer bg-[#00995D] text-white py-2 px-6 rounded-lg text-sm hover:bg-[#00814F] transition"
          on:click={handlePasswordChange}>
          Atualizar Senha
        </button>
      </div>
    </section>

    <!-- SEÇÃO SUPORTE (IVA WHATSAPP) -->
    <section class="shadow-lg rounded-2xl border border-gray-100 p-8 text-center" style="background-color: #E6F7EE;">
      <h2 class="text-2xl font-semibold text-gray-800 mb-4">Iva - Suporte via WhatsApp</h2>
      <p class="text-gray-600 mb-6">
        Converse com nossa assistente virtual no WhatsApp e receba ajuda imediata.
      </p>

      <img
        src={whatsapp}
        alt="WhatsApp"
        class="w-20 mx-auto mb-4"
      />

      <button
        class="cursor-pointer bg-[#00995D] text-white py-3 px-8 rounded-lg text-sm hover:bg-[#00814F] transition"
        on:click={openWhatsAppChat}>
        Abrir Chat no WhatsApp
      </button>
    </section>
  </div>
</main>
{/if}