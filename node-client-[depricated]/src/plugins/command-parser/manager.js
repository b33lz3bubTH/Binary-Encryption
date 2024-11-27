import { exec } from "child_process";
import os from "os";
import { getUniqueNameAndId } from '../../utils/sys.js'

export class CommandParser {
  constructor(telegramService, defaultChatId) {
    this.telegramService = telegramService;
    this.defaultChatId = defaultChatId;

    // Predefined commands
    this.commands = {
      hello: () => {
        this.sendTelegramMessage("Hello! This is your server!");
      },
      sysup: () => {
        exec("uptime", (error, stdout, stderr) => {
          this.handleExecResult(error, stdout, stderr);
        });
      },
      sysid: () => {
        const sysDetails = getUniqueNameAndId();
        this.sendTelegramMessage(`Username: ${sysDetails.username}\nSystem ID: ${sysDetails.systemId}\nDetails: ${JSON.stringify(sysDetails)}`);
      }
    };
  }

  async commandParser(cmd) {
    // Check if the command exists in predefined commands
    if (this.commands[cmd]) {
      try {
        await this.commands[cmd](); // Call the predefined function
      } catch (error) {
        this.sendTelegramMessage(`Error executing predefined command: ${error.message}`);
      }
    } else {
      // Fall back to executing shell command
      exec(cmd, (error, stdout, stderr) => {
        this.handleExecResult(error, stdout, stderr);
      });
    }
  }

  handleExecResult(error, stdout, stderr) {
    if (error) {
      this.sendTelegramMessage(`Error executing command: ${error.message}`);
      return;
    }
    if (stderr) {
      this.sendTelegramMessage(`Stderr: ${stderr}`);
      return;
    }
    this.sendTelegramMessage(`Stdout: ${stdout}`);
  }

  sendTelegramMessage(message) {
    this.telegramService.sendMessage(this.defaultChatId, message);
  }
}
