export function paymentRetryTimeout(attempt: number): number {
  return Math.min(1_000 * 2 ** attempt, 60_000);
}
