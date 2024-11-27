import os from 'os';
import crypto from "crypto";

export function getUniqueNameAndId() {
  try {
    // Get the username
    const username = os.userInfo().username;

    // Get system information to create a unique system ID
    const hostname = os.hostname();
    const platform = os.platform();
    const arch = os.arch();
    const cpus = os.cpus().map(cpu => cpu.model).join('-');
    const totalmem = os.totalmem();

    // Create a unique system ID using a combination of system information
    const systemId = [
      hostname,
      platform,
      arch,
      cpus,
      totalmem
    ].join('-').toLowerCase();

    const sys_unique_id = [
      hostname,
      platform,
    ].join('-').toLowerCase();

    const sys_md5_hash = crypto.createHash('md5').update(sys_unique_id).digest('hex');

    return {
      username,
      systemId,
      sys_md5_hash
    };
  } catch (error) {
    console.error('Error fetching system information:', error);
    throw error;
  }
}
