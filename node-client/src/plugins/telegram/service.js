import TelegramBot from 'node-telegram-bot-api';
import { appConfig } from '../../config/index.js';

export class TelegramNotifier {
  static instance;

  constructor() {
    this.bot = new TelegramBot(appConfig.TELEGRAM_BOT, { polling: false });
  }

  static getInstance() {
    if (!TelegramNotifier.instance) {
      TelegramNotifier.instance = new TelegramNotifier();
    }
    return TelegramNotifier.instance;
  }

  async sendMessage(chatId, message) {
    try {
      await this.bot.sendMessage(chatId, message);
      console.log('Message sent successfully');
    } catch (error) {
      console.error('Error sending message:', error.message);
    }
  }
}

