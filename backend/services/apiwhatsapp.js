"use strict";
const makeWASocket = require("@whiskeysockets/baileys").default;
const { useMultiFileAuthState, fetchLatestBaileysVersion, DisconnectReason } = require("@whiskeysockets/baileys");
const axios = require("axios");
const P = require("pino");
const qrcode = require("qrcode-terminal");

async function start() {
  const { state, saveCreds } = await useMultiFileAuthState("./baileys_auth");
  const { version } = await fetchLatestBaileysVersion();
  const sock = makeWASocket({
    version,
    auth: state,
    logger: P({ level: "silent" }),
  });

  sock.ev.on("creds.update", saveCreds);

  sock.ev.on("connection.update", ({ connection, lastDisconnect, qr }) => {
    if (qr) {
      console.log("Escaneie o QR abaixo para conectar:");
      qrcode.generate(qr, { small: true });
    }
    if (connection === "close") {
      const code = lastDisconnect?.error?.output?.statusCode || lastDisconnect?.error?.statusCode || 0;
      if (code !== DisconnectReason.loggedOut) {
        console.log("Conexão fechada, tentando reconectar...", code);
        start().catch((err) => console.error("Erro ao reconectar:", err));
      } else {
        console.log("Logout detectado. Remova a pasta 'baileys_auth' para reautenticar.");
      }
    } else if (connection === "open") {
      console.log("Conectado ao WhatsApp");
    } else if (connection === "connecting") {
      console.log("Conectando ao WhatsApp...");
    }
  });

  sock.ev.on("messages.upsert", async ({ messages, type }) => {
    if (type !== "notify") return;
    const msg = messages?.[0];
    if (!msg?.message) return;

    const remoteJid = msg.key?.remoteJid;
    if (!remoteJid || remoteJid === "status@broadcast") return;
    if (msg.key?.fromMe) return;

    // Extrai texto da mensagem
    let text = "";
    try {
      text =
        msg.message.conversation ||
        msg.message.extendedTextMessage?.text ||
        msg.message.imageMessage?.caption ||
        "";
    } catch (_) {}

    if (!text.trim()) return;

    const chatUrl = process.env.GO_CHAT_URL || "http://localhost:8080/chat";

    try {
      const resp = await axios.post(
        chatUrl,
        { message: text },
        { timeout: 30000, headers: { "Content-Type": "application/json" } }
      );
      const answer = resp.data?.response || "Sem resposta do serviço.";
      await sock.sendMessage(remoteJid, { text: answer }, { quoted: msg });
    } catch (err) {
      const code = err?.response?.status;
      const data = err?.response?.data;
      const errMsg = typeof data === "string" ? data : data?.error || err?.message || "erro desconhecido";
      console.error("Erro ao chamar serviço /chat:", code, errMsg);
      await sock.sendMessage(remoteJid, { text: `Erro ao consultar serviço: ${errMsg}` }, { quoted: msg });
    }
  });
}

start().catch((err) => console.error("Falha ao iniciar Baileys:", err));
