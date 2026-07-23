class Ledger:
    def __init__(self):
        self._entries = []

    @property
    def balance(self):
        return sum(entry["amount"] for entry in self._entries)

    def credit(self, amount):
        if amount <= 0:
            raise ValueError("credit amount must be positive")
        self._entries.append({"kind": "credit", "amount": amount})
