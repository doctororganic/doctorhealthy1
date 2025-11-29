type LogLevel = 'debug' | 'info' | 'warn' | 'error';

const prefix = '[doctorhealthy]';

function log(level: LogLevel, ...args: unknown[]) {
  // eslint-disable-next-line no-console
  const fn = level === 'error' ? console.error : level === 'warn' ? console.warn : level === 'info' ? console.info : console.debug;
  fn(prefix, ...args);
}

export const logger = {
  debug: (...args: unknown[]) => log('debug', ...args),
  info: (...args: unknown[]) => log('info', ...args),
  warn: (...args: unknown[]) => log('warn', ...args),
  error: (...args: unknown[]) => log('error', ...args),
};

// Compatibility namespace expected by some pages
export const loggers = {
  nutrition: logger,
  workouts: logger,
  plans: logger,
};
