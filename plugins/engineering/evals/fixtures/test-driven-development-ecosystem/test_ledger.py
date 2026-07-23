import unittest

from ledger import Ledger


class LedgerTests(unittest.TestCase):
    def test_credit_increases_balance(self):
        ledger = Ledger()
        ledger.credit(25)
        self.assertEqual(ledger.balance, 25)

    def test_credit_must_be_positive(self):
        ledger = Ledger()
        with self.assertRaisesRegex(ValueError, "credit amount must be positive"):
            ledger.credit(0)


if __name__ == "__main__":
    unittest.main()
