<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import Header from '$lib/components/layout/Header.svelte';
  import { browser } from '$app/environment';
  import iva from '$lib/assets/iva.png';

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

  let messages = [
    { from: 'ai', text: 'Olá! Eu sou a Iva (Intelbras Virtual Assistant). Como posso ajudá-lo hoje?' }
  ];
  let newMessage = '';
  let messagesContainer;

  let keepConnect = false;
  if (browser) {
    keepConnect = !!localStorage.getItem('cnt_email');
  }

  let inputDisabled = false; // controla se o input está bloqueado
  let showLimitWarning = false; // controla se mostramos a mensagem de aviso
  let aiTyping = false; // indica se a IA está "digitando"

  // Inicializa contador e bloqueia se já passou do limite
  const initSends = () => {
    const sends = parseInt(localStorage.getItem('input_sends') || '0', 10);
    if (sends >= 5 && !keepConnect) {
      inputDisabled = true;
      showLimitWarning = true;
    }
  };

  onMount(() => {
    initSends();
    scrollToBottom();
  });

  const sendMessage = async () => {
    // BLOQUEIA totalmente se estiver no limite
    if (inputDisabled || aiTyping) return;
    const question = newMessage.trim();
    if (!question) return;

    // Adiciona mensagem do usuário
    messages = [...messages, { from: 'user', text: question }];
    newMessage = '';
    scrollToBottom();

    // Atualiza contador no localStorage
    let sends = parseInt(localStorage.getItem('input_sends') || '0', 10);
    sends += 1;
    localStorage.setItem('input_sends', sends);

    // Se atingir limite, bloqueia input e mostra aviso
    if (sends >= 5 && !keepConnect) {
      inputDisabled = true;
      showLimitWarning = true;
    }

    // Mostra animação de digitando
    aiTyping = true;
    messages = [...messages, { from: 'ai', text: 'Digitando...', typing: true }];
    scrollToBottom();

    try {
      const email = browser ? (localStorage.getItem('cnt_email') || '') : '';
      const res = await fetch(`${API_BASE_URL}/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, message: question })
      });
      const data = await res.json().catch(() => null);
      console.log(data);
      console.log(res);
      // Remove indicador de digitando
      messages = messages.filter((m) => !m.typing);
      let answer = '';
      if (data) {
        if (data.status === 'success' && data.answer) answer = data.answer;
        else if (data.response) answer = data.response;
        else if (data.answer) answer = data.answer;
      }
      if (res.ok && answer) {
        messages = [...messages, { from: 'ai', text: answer }];
      } else {
        const msg = (data && (data.message || data.error)) ? (data.message || data.error) : 'Erro ao obter resposta da IA.';
        messages = [...messages, { from: 'ai', text: msg }];
      }
    } catch (err) {
      // Remove indicador de digitando e mostra erro
      messages = messages.filter((m) => !m.typing);
      messages = [...messages, { from: 'ai', text: 'Falha na conexão com o servidor.' }];
    } finally {
      aiTyping = false;
      scrollToBottom();
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') sendMessage();
  };

  const scrollToBottom = () => {
    setTimeout(() => {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }, 50);
  };

  const goToAccount = () => {
    goto('/account');
  };
</script>

<Header />

<main class="pt-28 pb-32 bg-gray-50 flex flex-col h-screen relative">

  <!-- Container das mensagens -->
  <div bind:this={messagesContainer} class="flex-1 overflow-y-auto p-6 space-y-4">
    {#each messages as msg}
      <div class="flex {msg.from === 'user' ? 'justify-end' : 'justify-start'}">
        {#if msg.from === 'ai'}
          <div class="flex items-end gap-2 max-w-[80%]">
            <img src={iva} alt="Iva" class="w-8 h-8 rounded-full border border-gray-200" />
            <div class="bg-white text-gray-800 p-3 rounded-lg shadow whitespace-pre-wrap max-w-md">
              {#if msg.typing}
                <div class="flex items-center space-x-1">
                  <span class="w-2 h-2 bg-gray-400 rounded-full animate-pulse"></span>
                  <span class="w-2 h-2 bg-gray-400 rounded-full animate-pulse" style="animation-delay: 0.2s"></span>
                  <span class="w-2 h-2 bg-gray-400 rounded-full animate-pulse" style="animation-delay: 0.4s"></span>
                </div>
              {:else}
                {msg.text}
              {/if}
            </div>
          </div>
        {:else}
          <div class="bg-[#00995D] text-white p-3 rounded-lg max-w-xs shadow whitespace-pre-wrap">
            {msg.text}
          </div>
        {/if}
      </div>
    {/each}
  </div>

  <!-- Aviso de limite acima do input -->
  {#if showLimitWarning}
    <button
      on:click={goToAccount}
      class="fixed bottom-20 left-1/2 transform -translate-x-1/2 cursor-pointer bg-yellow-100 border-l-4 border-yellow-400 text-yellow-800 p-3 rounded shadow max-w-md w-[90%] text-center z-30"
    >
      Você atingiu o limite de perguntas sem conta. Clique aqui para criar sua conta e ter acesso ilimitado.
    </button>
  {/if}

  <!-- Input de mensagem -->
  <div class="fixed bottom-0 left-0 w-full bg-white p-4 flex gap-2 items-center shadow-inner z-20">
    <input
      type="text"
      bind:value={newMessage}
      on:keypress={handleKeyPress}
      placeholder={inputDisabled ? "Crie sua conta para continuar enviando perguntas" : "Digite sua mensagem..."}
      class="flex-1 border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-[#00995D]"
      disabled={inputDisabled || aiTyping}
    />
    <button
      on:click={sendMessage}
      class="bg-[#00995D] hover:bg-[#00814F] text-white px-6 py-3 rounded-xl font-semibold transition-colors duration-300 disabled:opacity-50"
      disabled={inputDisabled || aiTyping}
    >
      Enviar
    </button>
  </div>
</main>
