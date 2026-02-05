# Decorator
- Decorator Design Pattern (mẫu Trang trí) là một mẫu thiết kế cấu trúc (Structural Pattern) trong OOP, dùng để mở rộng hành vi của một đối tượng một cách linh hoạt mà không cần sửa code gốc hay tạo quá nhiều subclass.

- Nói đơn giản:
👉 Bọc (wrap) một object bằng object khác để thêm chức năng cho nó, giống như mặc thêm áo khoác cho người vậy 🧥

## Ý tưởng cốt lõi

Giữ cùng interface với object gốc
Decorator chứa object gốc bên trong
Có thể xếp chồng nhiều decorator lên nhau

## Cấu trúc
- `Component` – interface hoặc abstract class
- `ConcreteComponent` – đối tượng gốc
- `Decorator` – abstract class, implement Component và chứa Component
- `ConcreteDecorator` – thêm hành vi mới

## Ví dụ dễ hiểu ☕ (cà phê)
- Giả sử ta có cà phê, và có thể thêm sữa, đường, kem…

```java

// Component
interface Coffee {
    String getDescription();
    double cost();
}

// ConcreteComponent
class BasicCoffee implements Coffee {
    public String getDescription() {
        return "Cà phê đen";
    }

    public double cost() {
        return 20000;
    }
}

// Decorator
abstract class CoffeeDecorator implements Coffee {
    protected Coffee coffee;

    public CoffeeDecorator(Coffee coffee) {
        this.coffee = coffee;
    }
}

// ConcreteDecorator
class MilkDecorator extends CoffeeDecorator {
    public MilkDecorator(Coffee coffee) {
        super(coffee);
    }

    public String getDescription() {
        return coffee.getDescription() + ", thêm sữa";
    }

    public double cost() {
        return coffee.cost() + 5000;
    }
}

// Sử dụng
Coffee coffee = new BasicCoffee();
coffee = new MilkDecorator(coffee);

System.out.println(coffee.getDescription()); // Cà phê đen, thêm sữa
System.out.println(coffee.cost());           // 25000

```
## Khi nào nên dùng Decorator?

✅ Khi:
Muốn thêm tính năng động cho object lúc runtime
Tránh tạo nhiều subclass kiểu CoffeeWithMilkAndSugarAndCream
Tuân thủ Open/Closed Principle (mở để mở rộng, đóng để sửa đổi)

❌ Không nên dùng khi:
Logic đơn giản, ít khả năng mở rộng
Decorator quá nhiều gây khó đọc code