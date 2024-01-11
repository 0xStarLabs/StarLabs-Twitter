
accs = []
with open("old.txt", "r") as f:
    accs = []
    for line in f:
        if 'AUTH_TOKEN: ' in line:
            accs.append(line.split('AUTH_TOKEN: ')[1])

with open("tokens.txt", "w") as f:
    for acc in accs:
        f.write(f"{acc}\n")