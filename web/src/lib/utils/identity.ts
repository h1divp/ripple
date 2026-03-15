const adjectives = ['Swift', 'Echoing', 'Silent', 'Vibrant', 'Mystic', 'Neon', 'Golden', 'Azure'];
const nouns = ['Runner', 'Whisper', 'Ghost', 'Pulse', 'Beacon', 'Shadow', 'Wave', 'Orbit'];

export function generateRandomName(): string {
  const adj = adjectives[Math.floor(Math.random() * adjectives.length)];
  const noun = nouns[Math.floor(Math.random() * nouns.length)];
  return `${adj} ${noun}`;
}

export function generateSessionSeed(): string {
  return crypto.randomUUID(); // for DiceBear seeds when unathenticated
}
