import { appConfig } from "./config/index.js";
import TelegramBot from 'node-telegram-bot-api';
import { TelegramNotifier } from "./plugins/telegram/service.js";
import { CommandParser } from "./plugins/command-parser/manager.js";

const bot = new TelegramBot(appConfig.TELEGRAM_BOT, { polling: true });
const telegramService = TelegramNotifier.getInstance();

const commandParser = new CommandParser(telegramService, appConfig.defaultChatId);

bot.on('message', async (msg) => {
  const chatId = msg.chat.id;
  if (chatId != appConfig.defaultChatId) {
    bot.sendMessage(chatId, 'You are not authorized to use this bot.');
  }
  console.log(`Chat ID: ${chatId}`);
  console.log(`Message: ${msg.text}`);
  commandParser.commandParser(msg.text.toLowerCase()).catch((err) => {
    console.log(`exec failed successfully. Error: ${err}`);
  });
});

