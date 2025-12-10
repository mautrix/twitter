package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000(\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0007\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u0005\n\u000b\f\r\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0003\u000e\u000f\u0010¨\u0006\u0011"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "Empty", "Leaf", "Unknown", "RatchetTreeLeafAdapter", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$Empty;", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$Leaf;", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class RatchetTreeLeaf implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new RatchetTreeLeafAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$Empty;", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "value", "Lcom/x/dmv2/thriftjava/EmptyNode;", "<init>", "(Lcom/x/dmv2/thriftjava/EmptyNode;)V", "getValue", "()Lcom/x/dmv2/thriftjava/EmptyNode;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Empty extends RatchetTreeLeaf {

        @InterfaceC88464a
        private final EmptyNode value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Empty(@InterfaceC88464a EmptyNode value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Empty copy$default(Empty empty, EmptyNode emptyNode, int i, Object obj) {
            if ((i & 1) != 0) {
                emptyNode = empty.value;
            }
            return empty.copy(emptyNode);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final EmptyNode getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Empty copy(@InterfaceC88464a EmptyNode value) {
            Intrinsics.m65272h(value, "value");
            return new Empty(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Empty) && Intrinsics.m65267c(this.value, ((Empty) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final EmptyNode m76807getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RatchetTreeLeaf(empty=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$Leaf;", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "value", "Lcom/x/dmv2/thriftjava/LeafNode;", "<init>", "(Lcom/x/dmv2/thriftjava/LeafNode;)V", "getValue", "()Lcom/x/dmv2/thriftjava/LeafNode;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Leaf extends RatchetTreeLeaf {

        @InterfaceC88464a
        private final LeafNode value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Leaf(@InterfaceC88464a LeafNode value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Leaf copy$default(Leaf leaf, LeafNode leafNode, int i, Object obj) {
            if ((i & 1) != 0) {
                leafNode = leaf.value;
            }
            return leaf.copy(leafNode);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final LeafNode getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Leaf copy(@InterfaceC88464a LeafNode value) {
            Intrinsics.m65272h(value, "value");
            return new Leaf(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Leaf) && Intrinsics.m65267c(this.value, ((Leaf) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final LeafNode m76808getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RatchetTreeLeaf(leaf=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$RatchetTreeLeafAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class RatchetTreeLeafAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public RatchetTreeLeaf m83775read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            RatchetTreeLeaf leaf;
            Intrinsics.m65272h(protocol, "protocol");
            RatchetTreeLeaf ratchetTreeLeaf = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    break;
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        ratchetTreeLeaf = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                    } else if (b == 12) {
                        leaf = new Leaf((LeafNode) LeafNode.ADAPTER.read(protocol));
                        ratchetTreeLeaf = leaf;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 12) {
                    leaf = new Empty((EmptyNode) EmptyNode.ADAPTER.read(protocol));
                    ratchetTreeLeaf = leaf;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
            if (ratchetTreeLeaf != null) {
                return ratchetTreeLeaf;
            }
            throw new IllegalStateException("unreadable");
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a RatchetTreeLeaf struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("RatchetTreeLeaf");
            if (struct instanceof Empty) {
                protocol.mo14136v3("empty", 1, (byte) 12);
                EmptyNode.ADAPTER.write(protocol, ((Empty) struct).m76807getValue());
            } else if (struct instanceof Leaf) {
                protocol.mo14136v3("leaf", 2, (byte) 12);
                LeafNode.ADAPTER.write(protocol, ((Leaf) struct).m76808getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTreeLeaf$Unknown;", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends RatchetTreeLeaf {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return 1210208356;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ RatchetTreeLeaf(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private RatchetTreeLeaf() {
    }
}