def main():
    with open("./books/frankenstein.txt") as f:
        file_contents = f.read()
        
        words = file_contents.split()
        print(len(words))

        lowered = file_contents.lower()
        print(lowered)

        char_count = {}
        for char in lowered:
            if char.isalpha():
                char_count[char] = char_count.get(char, 0) + 1
        
        sorted_pairs = sorted(char_count.items(), key=lambda kv: kv[1], reverse=True)
        for pair in sorted_pairs:
            print(f"The '{pair[0]}' character was found {pair[1]} times")
            
        def sort_on(dict):
            return dict["num"]
        
main()
