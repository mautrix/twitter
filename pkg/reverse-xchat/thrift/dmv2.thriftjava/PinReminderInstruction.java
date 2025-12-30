package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u00006\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000b\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0007\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\u0007\b\u0086\b\u0018\u0000 \u001c2\u00020\u0001:\u0002\u001d\u001cB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0005\u0010\u0006J\u0017\u0010\n\u001a\u00020\t2\u0006\u0010\b\u001a\u00020\u0007H\u0016¢\u0006\u0004\b\n\u0010\u000bJ\u0012\u0010\f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\rJ(\u0010\u000f\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u000f\u0010\u0010J\u0010\u0010\u0012\u001a\u00020\u0011HÖ\u0001¢\u0006\u0004\b\u0012\u0010\u0013J\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u0019\u001a\u00020\u00022\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u0019\u0010\u001aR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001bR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001b¨\u0006\u001e"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/PinReminderInstruction;", "Lcom/bendb/thrifty/a;", "", "should_register", "should_generate", "<init>", "(Ljava/lang/Boolean;Ljava/lang/Boolean;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Boolean;", "component2", "copy", "(Ljava/lang/Boolean;Ljava/lang/Boolean;)Lcom/x/dmv2/thriftjava/PinReminderInstruction;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Boolean;", "Companion", "PinReminderInstructionAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class PinReminderInstruction implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Boolean should_generate;

    @JvmField
    @InterfaceC88465b
    public final Boolean should_register;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new PinReminderInstructionAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/PinReminderInstruction$PinReminderInstructionAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/PinReminderInstruction;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/PinReminderInstruction;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/PinReminderInstruction;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class PinReminderInstructionAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public PinReminderInstruction m85667read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Boolean boolValueOf = null;
            Boolean boolValueOf2 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new PinReminderInstruction(boolValueOf, boolValueOf2);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 2) {
                        boolValueOf2 = Boolean.valueOf(protocol.readBool());
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 2) {
                    boolValueOf = Boolean.valueOf(protocol.readBool());
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a PinReminderInstruction struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("PinReminderInstruction");
            if (struct.should_register != null) {
                protocol.mo14136v3("should_register", 1, (byte) 2);
                protocol.mo14125P1(struct.should_register.booleanValue());
            }
            if (struct.should_generate != null) {
                protocol.mo14136v3("should_generate", 2, (byte) 2);
                protocol.mo14125P1(struct.should_generate.booleanValue());
            }
            protocol.mo14134i0();
        }
    }

    public PinReminderInstruction(@InterfaceC88465b Boolean bool, @InterfaceC88465b Boolean bool2) {
        this.should_register = bool;
        this.should_generate = bool2;
    }

    public static /* synthetic */ PinReminderInstruction copy$default(PinReminderInstruction pinReminderInstruction, Boolean bool, Boolean bool2, int i, Object obj) {
        if ((i & 1) != 0) {
            bool = pinReminderInstruction.should_register;
        }
        if ((i & 2) != 0) {
            bool2 = pinReminderInstruction.should_generate;
        }
        return pinReminderInstruction.copy(bool, bool2);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Boolean getShould_register() {
        return this.should_register;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Boolean getShould_generate() {
        return this.should_generate;
    }

    @InterfaceC88464a
    public final PinReminderInstruction copy(@InterfaceC88465b Boolean should_register, @InterfaceC88465b Boolean should_generate) {
        return new PinReminderInstruction(should_register, should_generate);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof PinReminderInstruction)) {
            return false;
        }
        PinReminderInstruction pinReminderInstruction = (PinReminderInstruction) other;
        return Intrinsics.m65267c(this.should_register, pinReminderInstruction.should_register) && Intrinsics.m65267c(this.should_generate, pinReminderInstruction.should_generate);
    }

    public int hashCode() {
        Boolean bool = this.should_register;
        int iHashCode = (bool == null ? 0 : bool.hashCode()) * 31;
        Boolean bool2 = this.should_generate;
        return iHashCode + (bool2 != null ? bool2.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "PinReminderInstruction(should_register=" + this.should_register + ", should_generate=" + this.should_generate + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
