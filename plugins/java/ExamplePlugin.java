import java.util.LinkedHashMap;
import java.util.Map;

public final class ExamplePlugin {
    public static void main(String[] args) {
        String target = args.length > 0 ? args[0] : "";
        System.out.println("{"plugin":"java-example","target":""
                + target.replace("\\","\\\\").replace(""","\\"")
                + "","type":"enrichment"}");
    }
}
